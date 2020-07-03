// Create by Yale 2018/3/1 9:41
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/yale8848/deploy/config"
	"github.com/yale8848/deploy/sshclient"
	zipfile "github.com/yale8848/deploy/zip"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const INSTALL_UNZIP_SHELL = `#!/usr/bin/env bash
function haveCmd() {
  if command -v $1 >/dev/null 2>&1;then
    return 1;
  else
    return 0;
  fi
}
function installUnzip () {
  haveCmd yum
  if [ $? == 1 ]
  then
    sudo yum -y install unzip
  else
    haveCmd apt-get
    if [ $? == 1 ]
    then
      sudo apt-get -y install unzip
    else
      echo "[deploy] #### please install unzip in server ####"
    fi
  fi

}
haveCmd unzip

if [ $? == 0 ]
then
 installUnzip
fi`

type verifyResult struct {
	Result  string
	Success bool
	Url     string
}
type comdOut struct {
	host    string
	hideMsg bool
}
type comdErr struct {
	host    string
	hideMsg bool
}

func (o comdOut) Write(p []byte) (n int, err error) {
	if !o.hideMsg {
		fmt.Print("[cmd] " + o.host + ": " + string(p))
	}
	return len(p), nil
}
func (o comdErr) Write(p []byte) (n int, err error) {
	if !o.hideMsg {
		fmt.Print("[cmd] " + o.host + ": " + string(p))
	}
	return len(p), nil
}
func checkError(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func getServers(c config.Config) []config.Server {
	cs := make([]config.Server, 0)
	for _, s := range c.Servers {
		ss := strings.Split(s.Host, ",")
		for _, sss := range ss {
			s.Host = sss
			cs = append(cs, s)
		}
	}
	return cs
}

func deleteFile(f string) {
	e := os.Remove(f)
	if e == nil {
		//fmt.Println("[delete] delete "+f+" success")
	}
}
func oneCmdMsgHide(cmd, host string, sc *sshclient.SSHClient) {

	if len(cmd) == 0 {
		return
	}
	out := comdOut{
		host:    host,
		hideMsg: true,
	}
	errw := comdErr{
		host:    host,
		hideMsg: true,
	}
	tcmd := cmd
	if !strings.HasPrefix(cmd, "sudo") {
		tcmd = "sudo " + cmd
	}
	err := sc.Command(tcmd, out, errw)
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(cmd+" success")
	}
}
func oneCmd(cmd, host string, sc *sshclient.SSHClient) {

	if len(cmd) == 0 {
		return
	}
	out := comdOut{
		host: host,
	}
	errw := comdErr{
		host: host,
	}

	tcmd := cmd
	if !strings.HasPrefix(cmd, "sudo") {
		tcmd = "sudo " + cmd
	}

	err := sc.Command(tcmd, out, errw)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tcmd + " success")
	}
}

func uploadOne(sc *sshclient.SSHClient, local, remote string, s config.Server) {
	uf := filepath.Base(local)
	err := sc.Upload(local, remote,
		func(percent int, finish bool) {
			fmt.Printf("[upload] %s: %s %d%%\r\n", s.Host, uf, percent)
			if finish {
				fmt.Printf("[upload] %s: %s finish\r\n", s.Host, uf)
			}
		})
	if err != nil {
		fmt.Println(err)
	}
	if remote[len(remote)-1] != '/' {
		remote = remote + "/"
	}
	cmd := fmt.Sprintf("unzip -qo %s -d %s", remote+uf, remote)
	oneCmd(cmd, s.Host, sc)
	cmd = fmt.Sprintf("rm -f %s", remote+uf)
	oneCmd(cmd, s.Host, sc)
	deleteFile(local)
}
func uploadSpecial(sc *sshclient.SSHClient, local, remote string, s config.Server) {
	uf := filepath.Base(local)
	sc.Upload(local, remote,
		func(percent int, finish bool) {
			//fmt.Printf("[upload] %s: %s %d%%\r\n",s.Host,uf,percent)
			if finish {
				//fmt.Printf("[upload] %s: %s finish\r\n",s.Host,uf)
			}
		})
	if remote[len(remote)-1] != '/' {
		remote = remote + "/"
	}
	oneCmdMsgHide("chmod 755 "+remote+uf, s.Host, sc)
	oneCmdMsgHide(remote+uf, s.Host, sc)
	cmd := fmt.Sprintf("rm -f %s", remote+uf)
	oneCmdMsgHide(cmd, s.Host, sc)
	deleteFile(local)

}
func installUnzip(sc *sshclient.SSHClient, s config.Server) {
	unzipShellName := os.Getenv("TMP") + "/installUnzip.sh"
	fout, err := os.Create(unzipShellName)
	if err != nil {
		return
	}
	fout.WriteString(INSTALL_UNZIP_SHELL)
	fout.Close()
	oneCmdMsgHide("mkdir /tmp", s.Host, sc)
	uploadSpecial(sc, unzipShellName, "/tmp", s)
}
func httpGet(url string) (string, bool) {
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	return string(body), true
}
func verify(s config.Server) verifyResult {
	result := verifyResult{}
	result.Url = s.Verify.Http + "://" + s.Host + s.Verify.Path
	if len(s.Verify.Path) == 0 {
		return result
	}
	if len(s.Verify.Http) == 0 {
		s.Verify.Http = "http"
	}
	if s.Verify.Delay == 0 {
		s.Verify.Delay = 3
	}
	if s.Verify.Count == 0 {
		s.Verify.Count = 3
	}
	if s.Verify.Gap == 0 {
		s.Verify.Gap = 2
	}
	time.Sleep(time.Duration(s.Verify.Delay) * time.Second)
	lastRet := ""
	for i := 0; i < s.Verify.Count; i++ {
		ret, success := httpGet(s.Verify.Http + "://" + s.Host + s.Verify.Path)
		if success {
			lastRet = ret
			if len(s.Verify.SuccessStrFlag) > 0 && strings.Contains(ret, s.Verify.SuccessStrFlag) {
				result.Success = true
				break

			}
		}
		time.Sleep(time.Duration(s.Verify.Gap) * time.Second)
	}
	result.Result = lastRet
	return result

}
func addJob(s config.Server, i int) {

	sc := sshclient.NewSSHClient()
	error := sc.ConnectTcp(s.Host, s.Port, s.User, s.Password, s.PrivateKeyPath, s.Socks5UrlPath, s.UserPasswordPath)
	checkError(error)

	installUnzip(sc, s)
	for _, cmd := range s.PreCommands {
		oneCmd(cmd, s.Host, sc)
	}

	for _, u := range s.Uploads {
		(func(up config.ServerUpload) {
			zipName := os.Getenv("TMP") + "/deploy_tmp_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ".zip"
			uploadOne(sc, zipfile.ZipFile(up.Local, zipName, up.ZipRegexp), up.Remote, s)
		})(u)
	}

	for _, cmd := range s.Commands {
		oneCmd(cmd, s.Host, sc)
	}
	sc.Close()
}

func main() {

	startTime := time.Now()
	configJson := flag.String("c", "config.json", "-c config.json")
	flag.Parse()
	configPath, _ := filepath.Abs(*configJson)
	config := config.Config{}
	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		checkError(err)
	}
	err = json.Unmarshal(dat, &config)
	if err != nil {
		checkError(err)
	}
	servers := getServers(config)

	var wg sync.WaitGroup

	if config.Concurrency {

		p, _ := ants.NewPool(len(servers))
		defer p.Release()

		for i, s := range servers {
			sss := s
			index := i
			wg.Add(1)

			p.Submit(func() {
				addJob(sss, index)
				wg.Done()
			})
		}

		wg.Wait()

	} else {
		for i, s := range servers {
			addJob(s, i)
		}
	}

	ct := time.Now().Sub(startTime).Seconds()
	fmt.Printf("[deploy] finish , total cost %d m : %d s\r\n", int(ct)/60, int(ct)%60)
	fmt.Println("start verify:")

	if config.Concurrency {

		verifyChan := make(chan verifyResult, len(servers))

		for _, s := range servers {
			sss := s
			go func() {
				verifyChan <- verify(sss)
			}()
		}
		for i := 0; i < len(servers); i++ {
			ret := <-verifyChan
			fmt.Printf("[verify] url: %s ,success: %t , result: %s \r\n", ret.Url, ret.Success, ret.Result)
		}
	} else {
		for _, s := range servers {
			ret := verify(s)
			fmt.Printf("[verify] url: %s ,success: %t , result: %s \r\n", ret.Url, ret.Success, ret.Result)
		}
	}

}

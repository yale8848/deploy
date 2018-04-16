// Create by Yale 2018/3/1 9:41
package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"flag"
	"path/filepath"
	"os"
	"time"
	"deploy/zip"
	"deploy/config"
	"deploy/sshclient"
	"strconv"
	"github.com/yale8848/gorpool"
)

type comdOut struct {
	host string
}
type comdErr struct {
	host string
}

func (o comdOut) Write(p []byte) (n int, err error)  {

	fmt.Print("[cmd] "+o.host+": "+string(p))
	return len(p),nil
}
func (o comdErr) Write(p []byte) (n int, err error)  {
	fmt.Print("[cmd] "+o.host+": "+string(p))
	return len(p),nil
}
func checkError(e error)  {
	if e !=nil{
		fmt.Println(e)
		panic(e)
	}
}

func  getServers(c config.Config) []config.Server {
	cs:= make([]config.Server,0)
	for _,s:=range c.Servers {
		ss:=strings.Split(s.Host,",")
		for _,sss:=range ss {
			s.Host = sss
			cs = append(cs,s)
		}
	}
	return cs
}

func zip(zipFiles []string,zipName string)  {

	if len(zipFiles) == 0 || len(zipName) == 0 {
		return
	}
	files:=[]*os.File{}
	for _,z:=range zipFiles  {
		f,e:=os.Open(z)
		if e!=nil {
			checkError(e)
		}
		files = append(files,f)
	}
	fmt.Println("[zip] start zip "+zipName)
	err:=zipfile.Compress(files,zipName)
	fmt.Println("[zip] zip "+zipName+" finish")
	checkError(err)
}

func deleteFile(f string)  {
	e:=os.Remove(f)
	if e==nil {
		fmt.Println("[delete] delete "+f+" success")
	}
}
func oneCmd(cmd,host string,sc *sshclient.SSHClient )  {

	if len(cmd)==0 {
		return
	}
	out:=comdOut{
		host:host,
	}
	errw:=comdErr{
		host:host,
	}

	err:=sc.Command(cmd,out,errw)
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println(cmd+" success")
	}
}
func addJob(s config.Server,i int)  {

	sc:=sshclient.NewSSHClient()
	error:=sc.ConnectTcp(s.Host,s.Port,s.User,s.Password)
	checkError(error)

	for _,cmd := range s.PreCommands{
		oneCmd(cmd,s.Host,sc)
	}

	for _,u:=range s.Uploads{
		(func(up config.ServerUpload) {
			t := time.Now()
			zipName := strconv.FormatInt(t.UTC().UnixNano(), 10)+".zip"
			zip(up.Local,zipName)
			uf:=filepath.Base(zipName)
			sc.Upload(zipName,up.Remote,
				func( percent int,finish bool) {
					fmt.Printf("[upload] %s: %s %d%%\r\n",s.Host,uf,percent)
					if finish {
						fmt.Printf("[upload] %s: %s finish\r\n",s.Host,uf)
					}
				})
			remote := up.Remote
			if remote[len(remote)-1] !='/'{
				remote =  remote+"/"
			}
			cmd:=fmt.Sprintf("unzip -qo %s -d %s",remote+zipName,remote)
			oneCmd(cmd,s.Host,sc)
			cmd=fmt.Sprintf("rm -f %s",remote+zipName)
			oneCmd(cmd,s.Host,sc)
			deleteFile(zipName)
		})(u)
	}

	for _,cmd := range s.Commands{
		oneCmd(cmd,s.Host,sc)
	}
	sc.Close()
}

func main()  {

	startTime:=time.Now()
	configJson:=flag.String("c","config.json","-c config.json")
	flag.Parse()
	configPath,_:=filepath.Abs(*configJson)
	config:=config.Config{}
	dat,err:=ioutil.ReadFile(configPath)
	if err!=nil {
		checkError(err)
	}
	err = json.Unmarshal(dat,&config)
	if err!=nil {
		checkError(err)
	}
	servers:=getServers(config)
	if config.Concurrency {
		pool:=gorpool.NewPool(len(servers),len(servers)).
			EnableWaitForAll(true).Start()
		for i,s:=range servers{
			sss:=s
			index:=i
			pool.AddJob(func() {
				addJob(sss,index)
			})
		}
		pool.WaitForAll()
	}else{
		for i,s:=range servers{
			addJob(s,i)
		}
	}


	ct:=time.Now().Sub(startTime).Seconds()
	fmt.Printf("finish: cost %d m : %d s",int(ct)/60,int(ct)%60)

}
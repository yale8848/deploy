// Create by Yale 2018/3/1 9:45
package sshclient

import (
	"errors"
	"github.com/pkg/sftp"
	"github.com/txthinking/socks5"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const maxMaxPacket = 1024 * 10

type SSHClient struct {
	client     *ssh.Client
	sftpClient *sftp.Client
}
type ClientError struct {
	msg string
}

type UploadCallback func(percent int, finish bool)

func (c ClientError) Error() string {
	return c.msg
}
func NewSSHClient() *SSHClient {
	return &SSHClient{}
}
func (s *SSHClient) Connect(network string, host string, port int, user string, password string,
	privateKeyPath string, socks5UrlPath string, userPasswordPath string) (err error) {

	if port == 0 {
		port = 22
	}

	if len(userPasswordPath) > 0 {
		by, err := ioutil.ReadFile(userPasswordPath)
		if err != nil {
			return err
		}
		str := strings.TrimSpace(string(by))

		ps := strings.Split(str, "@")
		if len(ps) > 0 {
			user = ps[0]
		}
		if len(ps) > 1 {
			password = ps[1]
		}
	}

	var client *ssh.Client

	cc := &ssh.ClientConfig{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		User: user,
	}

	if len(privateKeyPath) > 0 {

		pemBytes, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return err
		}

		cc.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		cc.Auth = []ssh.AuthMethod{ssh.Password(password)}
	}

	if len(socks5UrlPath) > 0 {

		sbs, err := ioutil.ReadFile(socks5UrlPath)
		if err != nil {
			return err
		}
		ur, err := url.Parse(string(sbs))
		if err != nil {
			return err
		}
		if ur.Scheme != "socks5" {
			return errors.New("scheme must socks5")
		}
		p, _ := ur.User.Password()
		tcpTimeout := 60
		tcpDeadline := 0
		udpDeadline := 60

		tt := ur.Query().Get("tcpTimeout")
		ttn, _ := strconv.ParseInt(tt, 10, 32)
		if ttn > 0 {
			tcpTimeout = int(ttn)
		}

		td := ur.Query().Get("tcpDeadline")
		tdn, _ := strconv.ParseInt(td, 10, 32)
		if tdn > 0 {
			tcpDeadline = int(tdn)
		}

		ud := ur.Query().Get("udpDeadline")
		udn, _ := strconv.ParseInt(ud, 10, 32)
		if udn > 0 {
			udpDeadline = int(udn)
		}

		c, _ := socks5.NewClient(ur.Host, ur.User.Username(), p, tcpTimeout, tcpDeadline, udpDeadline)
		conn, _ := c.Dial(network, host+":"+strconv.Itoa(port))
		cch, chans, reqs, err := ssh.NewClientConn(conn, network, cc)
		if err != nil {
			return err
		}
		client = ssh.NewClient(cch, chans, reqs)

	}

	if client == nil {
		client, err = ssh.Dial(network, host+":"+strconv.Itoa(port), cc)
	}

	if err == nil {
		s.client = client
		c, err := sftp.NewClient(client)
		if err == nil {
			s.sftpClient = c
			return nil
		} else {
			log.Fatal(err)
		}
	}
	return err

}
func (s *SSHClient) ConnectTcp(host string, port int, user string, password string, privateKeyPath string, socks5UrlPath string, userPasswordPath string) error {
	return s.Connect("tcp", host, port, user, password, privateKeyPath, socks5UrlPath, userPasswordPath)
}

func (s *SSHClient) Command(command string, out io.Writer, errW io.Writer) error {
	if s.client == nil {
		return ClientError{msg: "session not created"}
	}
	session, err := s.client.NewSession()
	if err == nil {
		session.Stdout = out
		session.Stderr = errW
		err = session.Run(command)
		if err == nil {
			session.Close()
		}
	}
	return err
}

func (s *SSHClient) Close() {
	s.client.Close()
	s.sftpClient.Close()
}

func (s *SSHClient) Upload(local string, remoteDir string, callback UploadCallback) error {
	srcFile, err := os.Open(local)
	if err != nil {
		return err
	}

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFile.Stat()

	remoteFileName := filepath.Base(local)
	dstFile, err := s.sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024*1024)

	var uploaded int64
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}

		uploaded += int64(n)

		dstFile.Write(buf[0:n])
		callback(int(uploaded*100/info.Size()), uploaded == info.Size())

	}
	return nil
}

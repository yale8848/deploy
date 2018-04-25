// Create by Yale 2018/3/1 9:45
package sshclient

import (
	"golang.org/x/crypto/ssh"
	"strconv"
	"os"
	"net"
	"github.com/pkg/sftp"
	"path"
	"path/filepath"
	"io"
)
const maxMaxPacket = (1 << 18) - 1024
type SSHClient struct {
	client *ssh.Client
	sftpClient *sftp.Client
}
type ClientError struct {
	msg string
}


type  UploadCallback func (percent int,finish bool)

func (c ClientError)Error()string  {
	return c.msg
}
func NewSSHClient()*SSHClient{
	return &SSHClient{}
}
func (s *SSHClient)Connect(network string,host string,port int,user string,password string) error {

	client, err := ssh.Dial(network, host+":"+strconv.Itoa(port), &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err==nil {
		s.client = client
		c, err := sftp.NewClient(client, sftp.MaxPacket(maxMaxPacket))
		if err == nil {
			s.sftpClient = c
			return nil
		}
	}
	return err

}
func (s *SSHClient)ConnectTcp(host string,port int,user string,password string) error {
	return s.Connect("tcp",host,port,user,password)
}

func (s *SSHClient)Command(command string,out io.Writer,errW io.Writer) error {
	if s.client==nil{
		return ClientError{msg:"session not created"}
	}
	session, err := s.client.NewSession()
	if err==nil {
		session.Stdout = out
		session.Stderr = errW
		err=session.Run(command)
		if err==nil {
			session.Close()
		}
	}
	return err
}

func (s *SSHClient)Close()  {
	s.client.Close()
	s.sftpClient.Close()
}

func (s *SSHClient)Upload(local string,remoteDir string,callback UploadCallback)error  {
	srcFile, err := os.Open(local)
	if err != nil {
		return err
	}

	info,err :=srcFile.Stat()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFile.Stat()

	remoteFileName :=filepath.Base(local)
	dstFile, err := s.sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err!=nil {
		return err
	}
	defer dstFile.Close()


	buf := make([]byte, 1024*1024)

	var uploaded int64
	for  {
		n, _ := srcFile.Read(buf)
		if n == 0{
			break
		}
		uploaded+=int64(n)
		dstFile.Write(buf[0:n])
		callback(int(uploaded*100/info.Size()),uploaded==info.Size())
	}
	return nil
}

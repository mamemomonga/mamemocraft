package actions

import (
	"golang.org/x/crypto/ssh"

	"bufio"
	"io/ioutil"
	"time"
	"log"
	"net"
	"io"
	"context"
	"fmt"
//	"github.com/davecgh/go-spew/spew"
)

type SSH struct {
	conf     *SSHConfig
	Client   *ssh.Client
	Session  *ssh.Session
}

type SSHConfig struct {
	KeyFile string
	User    string
	Host    string
	Port    string
	ConnectTimeout int
}

func NewSSH(c *SSHConfig) *SSH {
	t := new(SSH)
	t.conf = c
	return t
}

func (t *SSH)connect() (err error) {
	auth := []ssh.AuthMethod{}
	{
		key, err := ioutil.ReadFile(t.conf.KeyFile)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	sshConfig := &ssh.ClientConfig{
		User:            t.conf.User,
		Auth:            auth,
		Timeout:         time.Second * time.Duration(t.conf.ConnectTimeout),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", t.conf.Host + ":" + t.conf.Port, sshConfig)
	if err != nil {
		return err
	}
	t.Client = client
	return nil
}

func (t *SSH)session() (err error) {
	t.Session, err = t.Client.NewSession()
	if err != nil {
		return err
	}
	return nil
}

func (t *SSH)SessionClose() {
	t.Session.Close()
}

func (t *SSH)GetExitBool(cmd string) (bool,error) {
	err := t.connect()
	if err != nil {
		return false, err
	}
	defer t.Client.Close()

	err = t.session()
	if err != nil {
		return false, err
	}
	defer t.Session.Close()

	err = t.Session.Run(cmd)
	if err != nil {
		if ee,ok := err.(*ssh.ExitError); ok {
			if ee.Waitmsg.ExitStatus() == 1 {
				return true,nil
			}
		}
		return false,err
	}
	return false,nil
}

func (t *SSH)FileExists(path string) (bool,error) {
	ne,err := t.GetExitBool(fmt.Sprintf("test -e %s",path))
	if err != nil {
		return false, err
	}
	return !ne,nil
}

func (t *SSH)GetStdout(cmd string) (string, error) {
	err := t.connect()
	if err != nil {
		return "", err
	}
	defer t.Client.Close()

	err = t.session()
	if err != nil {
		return "", err
	}
	defer t.Session.Close()

	stdout, err := t.Session.StdoutPipe()
	if err != nil {
		return "",err
	}
	err = t.Session.Run(cmd)
	if err != nil {
		return "",err
	}
	s := bufio.NewScanner(stdout)

	buf:=""
	for s.Scan() {
		buf=buf+s.Text()+"\n"
	}
	return buf, nil
}

// https://stackoverflow.com/questions/21417223/simple-ssh-port-forward-in-golang
func (t *SSH)LocalForward(ctx context.Context, localAdr string, remoteAdr string) error {
	err := t.connect()
	if err != nil {
		return err
	}
	defer t.Client.Close()

	err = t.session()
	if err != nil {
		return err
	}
	defer t.Session.Close()

	localListener, err := net.Listen("tcp", localAdr)
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
		return err
	}
	defer localListener.Close()

	chDone := make(chan bool)

	forward := func(localConn net.Conn) {
		defer localConn.Close()

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Runtime Err: %v",err)
				chDone <- true
				return
			}
		}()

		sshConn, err := t.Client.Dial("tcp", remoteAdr)
		defer sshConn.Close()
		if err != nil {
			log.Printf("[t.Client.Dial]: %v", err)
			chDone <- true
			return
		}

		chS := make(chan bool)
		chR := make(chan bool)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[L->D] Runtime Err: %v",err)
					return
				}
			}()
			_, err = io.Copy(sshConn, localConn)
			if err != nil {
				log.Printf("[L->D] io.Copy failed: %v", err)
				return
			}
			chS <- true
		}()
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[D->L] Runtime Err: %v",err)
					return
				}
			}()
			_, err = io.Copy(localConn, sshConn)
			if err != nil {
				log.Printf("[D->L] io.Copy failed: %v", err)
				return
			}
			chR <- true
		}()
		<-chS // chSが帰ってくるまでブロッキング
		<-chR // chRが帰ってくるまでブロッキング
	}

	L:
    for {

		localConn, err := localListener.Accept()
		if err != nil {
			log.Fatalf("listen.Accept failed: %v", err)
		}

		go forward(localConn)

		select {
			case <-ctx.Done():
				break L
			case <-chDone:
				break L
			default:
		}

	}
	return nil
}

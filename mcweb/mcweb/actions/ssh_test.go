package actions_test

import (
	"os"
	//	"context"
	//	"time"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/actions"
	"log"
	"testing"
)

var sshconf *actions.SSHConfig

func TestSSH01(t *testing.T) {
	sshconf = &actions.SSHConfig{
		KeyFile:        "../../etc/id_ed25519",
		User:           "mamemocraft",
		Host:           "mc01.mamemo.online",
		Port:           "22",
		ConnectTimeout: 5,
	}
}
func TestSSH02(t *testing.T) {
	b, err := actions.NewSSH(sshconf).GetExitBool("exit 0")
	if err != nil {
		log.Printf("%#v", err)
		t.Error(err)
		os.Exit(1)
	}
	if b {
		log.Printf("True")
	} else {
		log.Printf("False")
	}
}

func TestSSH03(t *testing.T) {
	exist, err := actions.NewSSH(sshconf).FileExists("/home/mamemocraft/mamemocraft/var/running")
	if err != nil {
		log.Printf("%#v", err)
		t.Error(err)
		os.Exit(1)
	}
	if exist {
		log.Printf("File Exist")
	} else {
		log.Printf("File Not Exist")
	}
}

// func TestSSH02(t *testing.T) {
//
// 	ssh := actions.NewSSH(sshconf)
//
// 	exists,err := ssh.FileExists("/home/mamemocraft/mamemocraft/var/running")
//
// 	if err != nil {
// 		log.Println(err)
// 		t.Error(err)
// 		os.Exit(1)
// 	}
//
// 	if exists {
// 		t.Log("File Exists")
// 	} else {
// 		t.Log("File Not Exists")
// 	}
//
// }

// func TestSSHConnect(t *testing.T) {
// 	err := ssh.Connect()
// 	if err != nil {
// 		t.Error(err)
// 		os.Exit(1)
// 	}
// }
//
// func TestSSHSessionOpen(t *testing.T) {
// 	err := ssh.SessionOpen()
// 	if err != nil {
// 		t.Error(err)
// 		os.Exit(1)
// 	}
// }
//
// func TestSSHRun(t *testing.T) {
// 	err := ssh.Session.Run("test -e /home/mamemocraft/mamemocraft/var/running")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
//
// func TestSSHLocalForward(t *testing.T) {
//
// 	ctx := context.Background()
// 	ctx, cancel := context.WithTimeout(ctx, time.Second * 60)
// 	defer cancel()
//
// 	err := ssh.LocalForward(ctx,"127.0.0.1:5006","127.0.0.1:8123")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

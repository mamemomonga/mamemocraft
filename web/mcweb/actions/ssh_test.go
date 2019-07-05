package actions_test

import (
	"os"
	"testing"
	"context"
	"time"
	"github.com/mamemomonga/mamemocraft/web/mcweb/actions"
)

var ssh *actions.SSH

func TestSSHNewSSHConfig(t *testing.T) {
	ssh = actions.NewSSH(actions.NewSSHConfig{
		KeyFile: "../../etc/id_ed25519",
		User: "mamemocraft",
		Host: "mc01.mamemo.online",
		Port: "22",
		ConnectTimeout: 5,
	})
}

func TestSSHConnect(t *testing.T) {
	err := ssh.Connect()
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}
}

func TestSSHSessionOpen(t *testing.T) {
	err := ssh.SessionOpen()
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}
}

func TestSSHRun(t *testing.T) {
	err := ssh.Session.Run("test -e /home/mamemocraft/mamemocraft/var/running")
	if err != nil {
		t.Error(err)
	}
}

func TestSSHLocalForward(t *testing.T) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second * 60)
	defer cancel()

	err := ssh.LocalForward(ctx,"127.0.0.1:5006","127.0.0.1:8123")
	if err != nil {
		t.Error(err)
	}
}


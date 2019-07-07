package actions

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Sync struct {
	counter int
	ready   bool
	m       *sync.Mutex
	appdir  string
	enable  bool
}

type SyncConfig struct {
	Enable bool   `yaml:"enable"`
	APPDir string `yaml:"app_dir"`
}

func NewSync(c *SyncConfig) *Sync {
	t := new(Sync)
	t.enable = c.Enable
	t.appdir = c.APPDir
	t.m = new(sync.Mutex)
	return t
}

func (t *Sync) Start(ctx context.Context) {

	if !t.enable {
		return
	}

	time.Sleep(time.Second * 10)
	log.Printf("info: [Sync] Start")

	doit := func() {
		log.Printf("debug: [Sync] %d min.", t.counter)
		switch t.counter {
		case 10:
			log.Printf("info: [Sync] Run FIRST")
			t.runSync()

		case 70:
			log.Printf("info: [Sync] Run NORMAL")
			t.runSync()
			t.counter = 11
		}
	}
T:
	for {
		t.m.Lock()
		r := t.ready
		t.m.Unlock()
		if r {
			doit()
			t.counter++
		}
		time.Sleep(time.Minute)
		select {
		case <-ctx.Done():
			break T
		default:
		}
	}
	log.Printf("alert: [Sync] Terminate")
}

func (t *Sync) runSync() {
	err := t.runCommand(filepath.Join(t.appdir, "bin/sync.sh"))
	if err != nil {
		log.Printf("alert: [Sync] ERR %s", err)
	}
}

func (t *Sync) runCommand(c string, p ...string) error {
	cmd := exec.Command(c, p...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (t *Sync) Run() {
	if !t.ready {
		t.m.Lock()
		t.ready = true
		t.m.Unlock()
		t.counter = 0
	}
}

func (t *Sync) Stop() {
	t.m.Lock()
	t.ready = false
	t.m.Unlock()
}

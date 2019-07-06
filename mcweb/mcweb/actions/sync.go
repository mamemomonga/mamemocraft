package actions

import (
	"time"
	"log"
	"context"
	"sync"
	"os"
	"os/exec"
	"path/filepath"
)

type Sync struct {
	counter int
	ready   bool
	m       *sync.Mutex
	appdir  string
}

func NewSync(ad string)(*Sync) {
	t := new(Sync)
	t.appdir = ad
	t.m = new(sync.Mutex)
	return t
}

func (t *Sync) Start(ctx context.Context) {

	time.Sleep( time.Second * 10 )
	log.Printf("[SYNC] Start")

	doit := func() {
		log.Printf("[SYNC] %d",t.counter)
		switch t.counter {
		case 10:
		log.Printf("[SYNC] data file sync1")
		t.runSync()

		case 70:
		log.Printf("[SYNC] data file sync2")
		t.runSync()
		t.counter=11
		}
	}
	T:
	for {
		t.m.Lock()
		r := t.ready
		t.m.Unlock()
		if r {
			doit()
			t.counter ++
		}
		time.Sleep( time.Minute )
		select {
		case <-ctx.Done():
			break T
		default:
		}
	}
	log.Printf("[SYNC] Terminate")
}


func (t *Sync) runSync() {
	err := t.runCommand( filepath.Join(t.appdir,"bin/sync.sh"))
	if err != nil {
		log.Printf("[SYNC] ERR %s",err)
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
	if ! t.ready {
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

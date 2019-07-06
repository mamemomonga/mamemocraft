package actions

import (
	"time"
	"log"
	"context"
	"sync"
)

type Sync struct {
	counter int
	ready   bool
	m       *sync.Mutex
}

func NewSync()(*Sync) {
	t := new(Sync)
	t.m = new(sync.Mutex)
	log.Printf("[SYNC] Start")
	return t
}

func (t *Sync) Start(ctx context.Context) {
	doit := func() {
		log.Printf("[SYNC] %d",t.counter)
		switch t.counter {
		case 10:
		log.Printf("[SYNC] Ten!")
		case 70:
		log.Printf("[SYNC] Hour!")
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

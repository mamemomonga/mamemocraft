package actions

import (
	"context"
	"log"
	"net/http"
	"fmt"
)

const DymapModeNone = 0
const DymapModeWeb  = 1
const DymapModePF   = 2

type Dymap struct {
	conf    *DymapConfig
	mode    int
	cancel  context.CancelFunc
	doneCh  chan bool
}

type DymapConfig struct {
	Listen     string
	WebURL     string
	SSHconfig  *SSHConfig
}

func NewDymap(conf *DymapConfig)(*Dymap) {
	t := new(Dymap)
	t.conf = conf
	t.mode = DymapModeNone
	return t
}

func (t *Dymap) webHandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,`
<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<title>まめもくらふと</title>
<meta name="viewport" content="width=device-width">
</head>
<body>
<h1>Dymapはただいま利用できません</h1>	
<p><a href="%s">まめもくらふと</a></p>
</body>
</html>
`, t.conf.WebURL)
}


func (t *Dymap) RunWeb() {
	if t.mode == DymapModeWeb {
		return
	}

	if t.mode == DymapModePF {
		t.cancel() // PFの停止指示
		<-t.doneCh // PFの終了までブロック
	}
	t.doneCh = make(chan bool)
	t.mode   = DymapModeWeb

	mux := http.NewServeMux()
	mux.HandleFunc("/", t.webHandleIndex)
	srv := &http.Server{
		Addr: t.conf.Listen,
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		log.Printf("[dymap Web] START %s ",t.conf.Listen)
		defer func() {
			log.Println("[dymap Web] STOP")
			t.mode = DymapModeNone
			t.doneCh <- true
		}()
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("[dymap Web] %v",err)
		}
	}()

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
		t.mode = DymapModeNone
	}()
}

func (t *Dymap) RunPF() {
	if t.mode == DymapModePF {
		return
	}

	if t.mode == DymapModeWeb {
		t.cancel() // Webの停止指示
		<-t.doneCh // Webの停止までブロック
	}
	t.doneCh = make(chan bool)
	t.mode   = DymapModePF

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go func() {
		log.Printf("[dymap PF] START %s ", t.conf.Listen)
		defer func() {
			log.Println("[dymap PF] STOP")
			t.mode = DymapModeNone
			t.doneCh <- true
		}()
		ssh := NewSSH(t.conf.SSHconfig)
		err := ssh.Connect()
		if err != nil {
			log.Printf("[dymap PF] ssh.Connect %v",err)
			return
		}
		err = ssh.LocalForward(ctx, t.conf.Listen, "127.0.0.1:8123")
		if err != nil {
			log.Printf("[dymap PF] LocalForward %v",err)
		}
	}()
}


package actions

import (
	"fmt"
	"log"
	"time"
	"sync"
	"github.com/mamemomonga/mamemocraft/web/mcweb/web"
)

type Actions struct {
	gce      *GCE
	sshconf  *SSHConfig

	mutex    *sync.Mutex

	doStatus bool
	doStart  bool

	state    int
	message  string

	dymap    *Dymap
}

type Config struct {
	GCEKeyFile  string
	GCEProject  string
	GCEZone     string
	GCEInstance string
	SSHKeyFile  string
	SSHUser     string
	SSHHost     string
	SSHPort     string
}

const AutoReboot = true

const StatusUnknown  = 0
const StatusStop     = 1
const StatusLoading  = 2
const StatusRunning  = 3

func New(config Config) *Actions {
	t := new(Actions)

	t.gce = NewGCE()
	err := t.gce.LoadCredentialsFile(config.GCEKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	t.gce.Project  = config.GCEProject
	t.gce.Zone     = config.GCEZone
	t.gce.Instance = config.GCEInstance

	t.sshconf =  &SSHConfig{
		KeyFile: config.SSHKeyFile,
		User:    config.SSHUser,
		Host:    config.SSHHost,
		Port:    config.SSHPort,
		ConnectTimeout: 10,
	}

	t.dymap = NewDymap(&DymapConfig{
		Listen: "127.0.0.1:5006",
		WebURL: "https://mamemocraft-beta.mamemo.online/",
		SSHconfig: t.sshconf,
	})

	t.mutex = new(sync.Mutex)
	t.setStateMessage( StatusLoading, "準備中です。しばらくおまちください。")
	return t
}

func (t *Actions) setStateMessage(s int, m string) {
	t.mutex.Lock()
	t.state = s
	t.message = m
	t.mutex.Unlock()
}

func (t *Actions) Run() {
	t.doStatus = true

	go t.Runner()

	w := web.NewWebMain("127.0.0.1:5005")
	w.CbStatus = t.Status
	w.CbStart  = t.Start
	w.Run()
}

func (t *Actions) Runner() {

	for {
		log.Println("[RUNNER]")
		if t.doStatus {
			log.Println("[RUNNER] STATUS")
			t.chkStatus()
			t.mutex.Lock()
			t.doStatus = false
			t.mutex.Unlock()
		}
		if t.doStart {
			log.Println("[RUNNER] START")
			_, err := t.gce.Start()
			if err != nil {
				t.setStateMessage( StatusUnknown, "GCE インスタンス起動失敗😭")
				time.Sleep(time.Second * 60)
				return
			}
			t.setStateMessage( StatusLoading, "インスタンス起動中")
			t.mutex.Lock()
			t.doStart = false
			t.mutex.Unlock()
			time.Sleep(time.Second * 15)
		}

		if t.state == StatusRunning {
			t.dymap.RunPF()
		} else {
			t.dymap.RunWeb()
		}

		time.Sleep(time.Second * 10)
	}
}

func (t *Actions) chkStatus() {

	stateFile := func(name string) string {
		return fmt.Sprintf("/home/mamemocraft/mamemocraft/var/%s",name)
	}

	st, err := t.gce.Get()
	if err != nil {
		log.Printf("[GCE ERR] %s",err)
		return
	}
	log.Printf("[GCE] %s",st.Status)

	if st.Status == "RUNNING" {
		maintenance, err := t.sshFileExists(stateFile("maintenance"))
		if err != nil {
			return
		}
		if maintenance {
			t.setStateMessage( StatusUnknown, "ただいまメンテナンス中です")
			log.Printf("[SSH] mamemocraft is maintenance")
			return
		}
		stop, err := t.sshFileExists(stateFile("down"))
		if stop {
			t.setStateMessage( StatusLoading, "Minecraft Serverがとまってます😭")
			log.Printf("[SSH] mamemocraft is stop")
			if AutoReboot {
				_ = t.sshRun("sudo /sbin/reboot")
				log.Printf("[SSH] mamemocraft is rebooting")
				time.Sleep(time.Second * 30)
				t.setStateMessage( StatusLoading, "インスタンスを再起動しています")
				time.Sleep(time.Second * 30)
			}
			return
		}

		running, err := t.sshFileExists(stateFile("running"))
		if err != nil {
			return
		}
		if running {
			t.setStateMessage( StatusRunning, "まめもくらふと動作中")
			log.Printf("[SSH] mamemocraft is running")
			return
		} else {
			t.setStateMessage( StatusLoading, "Minecraft Serverを起動中")
			log.Printf("[SSH] mamemocraft is not running")
			return
		}
	}

	switch st.Status {
		case "STOPPING":
			t.setStateMessage( StatusUnknown, "停止作業中")
		case "TERMINATED":
			t.setStateMessage( StatusStop, "停止")
		default:
			t.setStateMessage( StatusUnknown, st.Status)
	}
	return
}

func (t *Actions) sshFileExists(path string) (exists bool, err error) {
	log.Printf("[SSH] ChkFile "+path)
	ssh := NewSSH(t.sshconf)
	err = ssh.Connect()
	if err != nil {
		log.Printf("[SSH] Connect %s",err)
		t.setStateMessage( StatusLoading, "状況わかんないです😭")
		return false,err
	}
	err = ssh.SessionOpen()
	if err != nil {
		log.Printf("[SSH] Session %s",err)
		t.setStateMessage( StatusLoading, "状況わかんないです😭")
		return false,err
	}
	defer ssh.SessionClose()
	err = ssh.Session.Run("test -e "+path)
	if err != nil {
		return false,nil
	}
	return true,nil
}

func (t *Actions) sshRun(cmd string) (err error) {
	log.Printf("[SSH] Run "+cmd)
	ssh := NewSSH(t.sshconf)
	err = ssh.Connect()
	err = ssh.Connect()
	if err != nil {
		log.Printf("[SSH] Connect %s",err)
	}
	err = ssh.SessionOpen()
	if err != nil {
		log.Printf("[SSH] Session %s",err)
	}
	defer ssh.SessionClose()
	err = ssh.Session.Run(cmd)
	log.Printf("[SSH] RetVal %s",err)
	return err
}

func (t *Actions) Status()(state int, message string) {
	state    = t.state
	message  = t.message
	t.mutex.Lock()
	t.doStatus = true
	t.mutex.Unlock()
	return
}

func (t *Actions) Start()(state int, message string) {
	t.setStateMessage( StatusLoading, "インスタンス起動")
	t.mutex.Lock()
	t.doStart = true
	t.mutex.Unlock()
	return
}


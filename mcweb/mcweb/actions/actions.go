package actions

import (
	"fmt"
	"log"
	"time"
	"sync"
	"context"
	"regexp"
	"strconv"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/web"
)

type Actions struct {
	gce      *GCE
	sshconf  *SSHConfig
	sync     *Sync

	mutex    *sync.Mutex

	doStatus bool
	doStart  bool

	mcRunning bool

	state    int
	message  string

	players  int
	playersZeroRemain int

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
const PlayersZeroEnable = true
const PlayersZeroMax = 10 // プレイヤーゼロが継続したらシャットダウン

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


	t.sync = NewSync()

	t.mcRunning = false
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
	t.playersZeroRemain = PlayersZeroMax

	go t.sync.Start(context.Background())
	go t.Runner()

	if PlayersZeroEnable {
		go t.playersZero()
	}

	w := web.NewWebMain("127.0.0.1:5005")
	w.CbStatus = t.Status
	w.CbStart  = t.Start
	w.Run()
}

func (t *Actions) Runner() {

	for {
		log.Println("[RUNNER]")
		if t.doStart {
			t.runnerDoStart()
		}
		if t.doStatus {
			t.runnerDoStatus()
		}
		if t.mcRunning {
			t.loginPlayers()
			go t.dymap.RunPF()
			t.sync.Run()
		} else {
			t.playersZeroRemain = PlayersZeroMax
			go t.dymap.RunWeb()
			t.sync.Stop()
		}
		time.Sleep(time.Second * 10)
	}
}

func (t *Actions) loginPlayers() {
	buf,err := NewSSH(t.sshconf).GetStdout("/home/mamemocraft/mamemocraft/bin/mcrcon -H localhost -p minecraft list")
	if err == nil {
		log.Printf("[Players] %s",buf)
		rex := regexp.MustCompile(`^There are (\d+) of a max (\d+) players online:`)
		match := rex.FindStringSubmatch(buf)
		current,_ := strconv.Atoi(match[1])
		t.players = current
		log.Printf("[Players] Current %d",t.players)
		if t.players > 0 {
			t.playersZeroRemain = PlayersZeroMax
		}
	}
}

func (t *Actions) playersZero() {
	for {
		if t.mcRunning {
			log.Printf("[PlayersZero] Remain: %d",t.playersZeroRemain)
			if t.playersZeroRemain == 0 {
				_,_ = NewSSH(t.sshconf).GetExitBool("sudo /sbin/poweroff")
			} else {
				t.playersZeroRemain--
			}
		}
		time.Sleep( time.Minute )
	}
}

func (t *Actions) runnerDoStart() {
	log.Println("[RUNNER] Start")
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
	time.Sleep(time.Second * 20)
}

func (t *Actions) runnerDoStatus() {
	log.Println("[RUNNER] Status")
	st, err := t.gce.Get()
	if err != nil {
		log.Printf("[GCE ERR] %s",err)
		t.setStateMessage( StatusLoading, "GCE 情報取得失敗😭")
		return
	}
	log.Printf("[GCE] %s",st.Status)

	switch st.Status {
	case "RUNNING":
		t.mcRunning =  t.mcStatus()
	case "STOPPING":
		t.mcRunning = false
		t.setStateMessage( StatusLoading, "停止作業中")
	case "TERMINATED":
		t.mcRunning = false
		t.setStateMessage( StatusStop, "停止")
	case "STAGING":
		t.mcRunning = false
		t.setStateMessage( StatusLoading, "起動準備中")
	default:
		t.mcRunning = false
		t.setStateMessage( StatusUnknown, st.Status)
	}
	t.mutex.Lock()
	t.doStatus = false
	t.mutex.Unlock()
}

func (t *Actions) mcStatus() bool {

	maintenance, err := t.sshFileExists("maintenance")
	if err != nil {
		return false
	}
	if maintenance {
		t.setStateMessage( StatusUnknown, "ただいまメンテナンス中です")
		log.Printf("[SSH] mamemocraft is maintenance")
		return false
	}
	stop, err := t.sshFileExists("down")
	if stop {
		t.setStateMessage( StatusLoading, "Minecraft Serverがとまってます😭")
		log.Printf("[SSH] mamemocraft is stop")
		if AutoReboot {
			_,_ = NewSSH(t.sshconf).GetExitBool("sudo /sbin/reboot")
			log.Printf("[SSH] mamemocraft is rebooting")
			time.Sleep(time.Second * 30)
			t.setStateMessage( StatusLoading, "インスタンスを再起動しています")
			time.Sleep(time.Second * 30)
		}
		return false
	}

	running, err := t.sshFileExists("running")
	if err != nil {
		return false
	}
	if running {
		t.setStateMessage( StatusRunning, "まめもくらふと動作中")
		log.Printf("[SSH] mamemocraft is running")
		return true
	} else {
		t.setStateMessage( StatusLoading, "Minecraft Serverを起動中")
		log.Printf("[SSH] mamemocraft is not running")
		return false
	}
}

func (t *Actions) sshFileExists(name string) (exists bool, err error) {
	log.Printf("[SSH] ChkFile %s",name)
	path:=fmt.Sprintf("/home/mamemocraft/mamemocraft/var/%s",name)

	exists, err = NewSSH(t.sshconf).FileExists(path)

	if err != nil {
		log.Printf("[SSH] Error %s",err)
		t.setStateMessage( StatusLoading, "状況わかんないです😭")
		return false, err
	}
	return exists,nil
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


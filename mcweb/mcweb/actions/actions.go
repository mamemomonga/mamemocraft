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
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/config"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/buildinfo"
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

	rconPassword string

	dymap    *Dymap
	mastodon *Mastodon
	prevToot string

}

const AutoReboot = true
const PlayersZeroEnable = true
const PlayersZeroMax = 10 // プレイヤーゼロが継続したらシャットダウン

const StatusUnknown  = 0
const StatusStop     = 1
const StatusLoading  = 2
const StatusRunning  = 3

func New(configFile string) *Actions {
	t := new(Actions)

	config,err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}

	t.rconPassword = config.Rcon.Password

	t.gce = NewGCE()
	err = t.gce.LoadCredentialsFile(config.GCE.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	t.gce.Project  = config.GCE.Project
	t.gce.Zone     = config.GCE.Zone
	t.gce.Instance = config.GCE.Instance

	t.sshconf =  &SSHConfig{
		KeyFile: config.SSH.KeyFile,
		User:    config.SSH.User,
		Host:    config.SSH.Host,
		Port:    config.SSH.Port,
		ConnectTimeout: 10,
	}

	t.dymap = NewDymap(&DymapConfig{
		Listen: config.Dymap.Listen,
		WebURL: config.Dymap.WebURL,
		SSHconfig: t.sshconf,
	})

	t.sync = NewSync(&SyncConfig{
		Enable: config.Sync.Enable,
		APPDir: config.Sync.APPDir,
	})

	t.mastodon = NewMastodon( &MastodonConfig{
		Server:     config.Mastodon.Server,
		Email:      config.Mastodon.Email,
		Password:   config.Mastodon.Password,
		ClientFile: config.Mastodon.ClientFile,
		ClientName: config.Mastodon.ClientName,
	})

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

func (t *Actions) toot(s string) {

	if t.prevToot == s {
		return
	}

	if err := t.mastodon.Toot( fmt.Sprintf("[まめもくらふと] %s ﾖｼ :genbaneko:",s)); err != nil {
		log.Printf("alert: [Mastodon] %s", err)
	}

	t.prevToot = s
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

	if err := t.mastodon.Connect(); err != nil {
		log.Printf("alert: [Mastodon] %s",err)
	}

	w.TPData = map[string]string {
		"AppName": "mamemocraft-web",
		"Version": buildinfo.Version,
		"Revision": buildinfo.Revision,
	}
	w.Run() // ブロック
}

func (t *Actions) Runner() {

	for {
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
	cmd := fmt.Sprintf("/home/mamemocraft/mamemocraft/bin/mcrcon -H localhost -p %s list",t.rconPassword)
	buf,err := NewSSH(t.sshconf).GetStdout(cmd)
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
		t.toot("停止")

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
		t.toot("起動")
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


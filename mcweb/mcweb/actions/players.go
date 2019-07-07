package actions

import (
	"log"
	"regexp"
	"strconv"
	"time"
	"strings"
	"github.com/thoas/go-funk"
)

type Players struct {
	c      *PlayersConfig
	Remain int
	Count  int
	prevPlayers []string
}

type PlayersConfig struct {
	ZeroShutdown      bool
	ZeroShutdownCount int
	MCRunning         *bool
	Shutdown          func()
	Toot              func(s string)
}

func NewPlayers(c *PlayersConfig) *Players {
	t := new(Players)
	t.c = c
	return t
}

func (t *Players) Reset() {
	t.Remain = t.c.ZeroShutdownCount
}

func (t *Players) Start() {
	t.Reset()
	for {
		if *t.c.MCRunning {
			log.Printf("[PlayersZero] Remain: %d", t.Remain)
			if t.Remain == 0 {
				t.c.Shutdown()
			} else {
				t.Remain--
			}
		}
		time.Sleep(time.Minute)
	}
}

func (t *Players)Check(buf string) {
	log.Printf("[Players] %s",buf)
	rex := regexp.MustCompile(`^There are (\d+) of a max (\d+) players online:(.+)`)

	match := rex.FindStringSubmatch(buf)
	current,_ := strconv.Atoi(match[1])
	t.Count = current

	log.Printf("[Players] Current %d", t.Count)
	if t.Count > 0 {
		t.Reset()
	}

	names := funk.Map(strings.Split(match[3],","),func(x string) string {
		return strings.TrimSpace(x)
	}).([]string)

	names = funk.FilterString(names,func(x string) bool {
		if x == "" {
			return false
		}
		return true
	})

	pnew   := funk.FilterString(names,func(x string) bool { return ! funk.ContainsString(t.prevPlayers,x) })
	pleave := funk.FilterString(t.prevPlayers,func(x string) bool { return ! funk.ContainsString(names,x) })

	if len(pnew) > 0 {
		log.Printf("[Players] New Member: %#v", pnew)
		t.c.Toot( t.san(pnew)+" 参加")
	}
	if len(pleave) > 0 {
		log.Printf("[Players] Leave Member: %#v", pleave)
		t.c.Toot( t.san(pleave)+" 退出")
	}
	t.prevPlayers = names
}

func (t *Players)san(m []string) string {

	return strings.Join(funk.Map(m,func(x string) string {
		return x+"さん"
	}).([]string),",")
}


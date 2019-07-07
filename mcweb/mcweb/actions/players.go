package actions

import(
	"regexp"
	"strconv"
	"log"
	"time"
)

type Players struct {
	c *PlayersConfig
	Remain int
	Count  int
}

type PlayersConfig struct {
	ZeroShutdown      bool
	ZeroShutdownCount int
	MCRunning         *bool
	Shutdown          func()
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
			log.Printf("[PlayersZero] Remain: %d",t.Remain)
			if t.Remain == 0 {
				t.c.Shutdown()
			} else {
				t.Remain--
			}
		} else {
			log.Printf("[PlayersZero] MC Not Running")
		}

		time.Sleep( time.Minute )
	}
}

func (t *Players) Check(buf string) {
	log.Printf("[Players] %s",buf)
	rex := regexp.MustCompile(`^There are (\d+) of a max (\d+) players online:`)
	match := rex.FindStringSubmatch(buf)
	current,_ := strconv.Atoi(match[1])
	t.Count = current

	log.Printf("[Players] Current %d",t.Count)
	if t.Count > 0 {
		t.Reset()
	}
}


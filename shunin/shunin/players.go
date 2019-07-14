package main

import (
	"log"
	"regexp"
	"strconv"
	"time"
	"strings"
	"github.com/thoas/go-funk"
	"github.com/mamemomonga/mamemocraft/shunin/shunin/rcon"
)

type Players struct {
	c      *PlayersConfig
	prevPlayers []string
	Count  int
}

type PlayersConfig struct {
	Toot func(s string)
}

func NewPlayers(c *PlayersConfig) *Players {
	t := new(Players)
	t.c = c
	return t
}

func (t *Players) Run() {
	for {
		t.Check()
		time.Sleep(time.Second * 10)
	}
}

func (t *Players) Check() {

	list, err := rcon.Get("list")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Players] %s",list)
	rex := regexp.MustCompile(`^There are (\d+) of a max (\d+) players online:(.+)`)

	match := rex.FindStringSubmatch(list)
	current,_ := strconv.Atoi(match[1])
	t.Count = current

	log.Printf("[Players] Current %d", t.Count)

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


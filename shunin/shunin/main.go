package main

import (
	cs "github.com/mamemomonga/mamemocraft/shunin/shunin/confstate"
	"github.com/mamemomonga/mamemocraft/shunin/shunin/rcon"
	"github.com/mamemomonga/mamemocraft/shunin/shunin/mastodon"
	"log"
	"fmt"
)

func main() {
	if err := cs.Load(""); err != nil {
		log.Fatal(err)
	}

	rcon.Host = cs.C().Rcon.Host
	rcon.Password = cs.C().Rcon.Password
	if c,err := cs.GetDir("mcrcon"); err == nil {
		rcon.Mcrcon = c
	} else {
		log.Fatal(err)
	}

	var mastodonClientFile string
	if c,err := cs.GetDir("etc/shunin.json"); err == nil {
		mastodonClientFile = c
	} else {
		log.Fatal(err)
	}

	m := mastodon.NewMastodon(&mastodon.MastodonConfig{
		Server: cs.C().Shunin.Server,
		Email:  cs.C().Shunin.Email,
		Password: cs.C().Shunin.Password,
		ClientName: "まめもくらふと",
		ClientFile: mastodonClientFile,
	})

	if err := m.Connect(); err != nil {
		log.Fatal(err)
	}

	toot := func(s string) {
		if err := m.Toot(fmt.Sprintf("[まめもくらふと]\n%s ﾖｼ!:genbaneko:",s)); err != nil {
			log.Println(err)
		}
	}

	p := NewPlayers(&PlayersConfig{
		Toot: toot,
	})

	p.Run()
}


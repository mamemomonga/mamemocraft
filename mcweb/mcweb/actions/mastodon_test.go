package actions_test

import (
	"log"
	"testing"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/actions"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/config"
)

// go test -v --count=1 mcweb/actions/mastodon_test.go

var don *actions.Mastodon

func init() {
	c,err := config.Load("../../etc/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mastodon_config := actions.MastodonConfig{
		Server:     c.Mastodon.Server,
		Email:      c.Mastodon.Email,
		Password:   c.Mastodon.Password,
		ClientFile: "../../etc/mastodon.json",
		ClientName: c.Mastodon.ClientName,
	}
	don = actions.NewMastodon( &mastodon_config )
	if err := don.Connect(); err != nil {
		log.Fatal(err)
	}

}

func TestMastodon01(t *testing.T) {
	don.Toot("テスト")
}


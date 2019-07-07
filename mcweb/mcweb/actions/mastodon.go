package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-mastodon"
	"io/ioutil"
	"log"
	"os"
)

type Mastodon struct {
	c        *MastodonConfig
	client   *mastodon.Client
	Ready    bool
	lastToot string
}

type MastodonConfig struct {
	Enable     bool
	Server     string
	Email      string
	Password   string
	ClientFile string
	ClientName string
}

type ClientConfigs struct {
	Tokens map[string]ClientTokens `json:"tokens"`
}

type ClientTokens struct {
	ClientID     string `json:"id"`
	ClientSecret string `json:"secret"`
}

func NewMastodon(c *MastodonConfig) *Mastodon {
	t := new(Mastodon)
	t.c = c
	t.Ready = false
	return t
}

func (t *Mastodon) Connect() (err error) {
	if !t.c.Enable {
		return nil
	}

	ctx := context.Background()

	ccs := &ClientConfigs{
		Tokens: make(map[string]ClientTokens),
	}

	if _, err := os.Stat(t.c.ClientFile); !os.IsNotExist(err) {
		if err := t.loadClientFile(ccs); err != nil {
			return err
		}
	}
	if _, ok := ccs.Tokens[t.c.Server]; !ok {
		app, err := mastodon.RegisterApp(ctx, &mastodon.AppConfig{
			Server:     fmt.Sprintf("https://%s/", t.c.Server),
			ClientName: t.c.ClientName,
			Scopes:     "read write follow",
		})
		if err != nil {
			return err
		}
		ccs.Tokens[t.c.Server] = ClientTokens{
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
		}
		if err := t.saveClientFile(ccs); err != nil {
			return err
		}
	}
	t.client = mastodon.NewClient(&mastodon.Config{
		Server:       fmt.Sprintf("https://%s/", t.c.Server),
		ClientID:     ccs.Tokens[t.c.Server].ClientID,
		ClientSecret: ccs.Tokens[t.c.Server].ClientSecret,
	})
	if err := t.client.Authenticate(ctx, t.c.Email, t.c.Password); err != nil {
		return err
	}
	account, err := t.client.GetAccountCurrentUser(ctx)
	if err != nil {
		return err
	}
	log.Printf("info: [Mastodon] Server:      %s", t.c.Server)
	log.Printf("info: [Mastodon] SelfID:      %s", account.ID)
	log.Printf("info: [Mastodon] Username:    %s", account.Username)
	log.Printf("info: [Mastodon] DisplayName: %s", account.DisplayName)
	t.Ready = true
	return nil
}

func (t *Mastodon) saveClientFile(cc *ClientConfigs) (err error) {
	buf, err := json.Marshal(cc)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(t.c.ClientFile, buf, 0644)
	if err != nil {
		return
	}
	log.Printf("info: [Mastodon] Save ClientFile")
	return nil
}

func (t *Mastodon) loadClientFile(cc *ClientConfigs) (err error) {
	buf, err := ioutil.ReadFile(t.c.ClientFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, cc)
	if err != nil {
		return
	}
	log.Printf("info: [Mastodon] Load ClientFile")
	return nil
}

func (t *Mastodon) TootNoDup(s string) error {
	if t.lastToot == s {
		return nil
	}
	if err := t.Toot(s); err != nil {
		return err
	}
	t.lastToot = s
	return nil
}

func (t *Mastodon) Toot(s string) error {
	if !t.c.Enable {
		return nil
	}
	if !t.Ready {
		return nil
	}
	ctx := context.Background()
	toot := mastodon.Toot{Status: s}
	_, err := t.client.PostStatus(ctx, &toot)
	if err != nil {
		return err
	}
	log.Printf("info: [Mastodon] Say: %s", toot.Status)
	return nil
}

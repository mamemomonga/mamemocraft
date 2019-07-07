package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	GCE        GCEType      `yaml:"gce"`
	SSH        SSHType      `yaml:"ssh"`
	Sync       SyncType     `yaml:"sync"`
	Rcon       RconType     `yaml:"rcon"`
	Dymap      DymapType    `yaml:"dymap"`
	Mastodon   MastodonType `yaml:"mastodon"`
	Players    PlayersType  `yaml:"players"`
	AutoReboot bool         `yaml:"auto_reboot"`
}

type GCEType struct {
	KeyFile   string `yaml:"keyfile"`
	Project   string `yaml:"project"`
	Zone      string `yaml:"zone"`
	Instance  string `yaml:"instance"`
}

type SSHType struct {
	KeyFile   string `yaml:"keyfile"`
	User      string `yaml:"user"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
}

type SyncType struct {
	Enable   bool   `yaml:"enable"`
	APPDir   string `yaml:"app_dir"`
}

type RconType struct {
	Password string `yaml:"password"`
}

type DymapType struct {
	Listen     string `yaml:"listen"`
	Forward    string `yaml:"forward"`
	WebURL     string `yaml:"weburl"`
}

type MastodonType struct {
	Enable      bool   `yaml:"enable"`
	Server      string `yaml:"server"`
	Email       string `yaml:"email"`
	Password    string `yaml:"password"`
	ClientFile  string `yaml:"client_file"`
	ClientName  string `yaml:"client_name"`
}

type PlayersType struct {
	ZeroShutdown bool `yaml:"zero_shutdown"`
	ZeroShutdownCount int `yaml:"zero_shutdown_count"`
}


func Load(filename string) (data *Config, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return
	}
	log.Printf("Read: %s", filename)
	return data, nil
}


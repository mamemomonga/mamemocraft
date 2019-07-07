package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	GCE   GCEType   `yaml:"gce"`
	SSH   SSHType   `yaml:"ssh"`
	Sync  SyncType  `yaml:"sync"`
	Rcon  RconType  `yaml:"rcon"`
	Dymap DymapType `yaml:"dymap"`
}

type GCEType struct {
	KeyFile   string `yaml:"key_file"`
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


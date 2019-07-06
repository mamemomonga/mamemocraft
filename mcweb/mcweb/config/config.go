package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	GCEKeyFile   string `yaml:"gce_key_file"`
	GCEProject   string `yaml:"gce_project"`
	GCEZone      string `yaml:"gce_zone"`
	GCEInstance  string `yaml:"gce_instance"`
	SSHKeyFile   string `yaml:"ssh_keyfile"`
	SSHUser      string `yaml:"ssh_user"`
	SSHHost      string `yaml:"ssh_host"`
	SSHPort      string `yaml:"ssh_port"`
	SyncAPPDir   string `yaml:"sync_app_dir"`
	RConPassword string `yaml:"rcon_password"`
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

package main

import (
	"log"

	"github.com/mamemomonga/mamemocraft/web/mcweb/buildinfo"
	"github.com/mamemomonga/mamemocraft/web/mcweb/web"
	"github.com/mamemomonga/mamemocraft/web/mcweb/actions"
)

func main() {

	log.Printf("mamemocraft-web Version: %s Revision: %s\n", buildinfo.Version, buildinfo.Revision)

	act := actions.New(actions.Config{
		GCEKeyFile:  "./etc/gce-key.json",
		GCEProject:  "mamemo-190623",
		GCEZone:     "asia-northeast1-b",
		GCEInstance: "mamemocraft-190624",
		SSHKeyFile:  "./etc/id_ed25519",
		SSHUser:     "mamemocraft",
		SSHHost:     "mc01.mamemo.online",
		SSHPort:     "22",
	})
	act.Run()

//	go func() {
//		wdy := web.NewWebDymap("0.0.0.0:5006")
//		wdy.Run()
//	}()

	web.NewWebMain("0.0.0.0:5005",act).Run()
}


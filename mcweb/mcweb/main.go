package main

import (
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/actions"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/buildinfo"
	"log"
	"os"
)

func main() {
	log.Printf("mamemocraft-web Version: %s Revision: %s\n", buildinfo.Version, buildinfo.Revision)

	if len(os.Args) != 2 {
		log.Printf("USAGE: mcweb config.yaml\n")
		os.Exit(1)
	}
	actions.New(os.Args[1]).Run()

	//	act := actions.New(actions.Config{
	//		GCEKeyFile:  "./etc/gce-key.json",
	//		GCEProject:  "mamemo-190623",
	//		GCEZone:     "asia-northeast1-b",
	//		GCEInstance: "mamemocraft-190624",
	//		SSHKeyFile:  "./etc/id_ed25519",
	//		SSHUser:     "mamemocraft",
	//		SSHHost:     "mc01.mamemo.online",
	//		SSHPort:     "22",
	//		SyncAPPDir:  "../sync",
	//	})
	//	act.Run()

}

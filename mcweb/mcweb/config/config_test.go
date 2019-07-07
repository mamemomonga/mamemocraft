package config_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mamemomonga/mamemocraft/mcweb/mcweb/config"
	"log"
	"testing"
)

// go test -v --count=1 mcweb/config/config_test.go

var cnf *config.Config

func init() {
	c, err := config.Load("../../etc/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	cnf = c
}

func TestConfig01(t *testing.T) {
	spew.Dump(cnf)
}

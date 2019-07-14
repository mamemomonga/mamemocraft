package rcon

import (
	"os/exec"
	"bufio"
)

var (
	Host     string
	Password string
	Mcrcon   string
)

func Get(c string) (string,error) {

	cmd := exec.Command(Mcrcon, "-H", Host, "-p", Password,"-c",c)

	stdout,err := cmd.StdoutPipe()
	if err != nil {
		return "",err
	}
	if err := cmd.Start(); err != nil {
		return "",err
	}

	b:=""
	s:=bufio.NewScanner(stdout)
	for s.Scan() {
		b=b+s.Text()
	}
	return b,nil
}


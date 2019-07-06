package web

import (
	packr "github.com/gobuffalo/packr/v2"
	"log"
)

func boxStatic() *packr.Box {
	return packr.New("static","../../assets/static")
}

func boxTemplates(find string) (string, error) {
	b := packr.New("templates","../../assets/templates")
	s, err := b.FindString(find)
	if err != nil {
		log.Printf("warn: [WebBox] %s",err)
	}
	return s,err
}


package web

import (
	"log"
	"fmt"
	"net/http"
	"html/template"
)

type WebDymap struct {
	Server *http.Server
}

func NewWebDymap(listen string) *WebDymap {
	t := new(WebDymap)

	mux := http.NewServeMux()
	mux.HandleFunc("/",t.handleIndex)

	t.Server = &http.Server{
		Addr: listen,
		Handler: mux,
	}

	return t
}

func (t *WebDymap) Run() {
	log.Printf("[WebDymap] Start Listening at %s", t.Server.Addr)
	err := t.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (t *WebDymap) handleIndex(w http.ResponseWriter, r *http.Request) {
	s, err := boxTemplates("index_dymap.tpl.html")
	if err != nil {
		fmt.Fprintf(w,"[WebDymap] Template error %v",err)
	}

	tp, err := template.New("T").Parse(s)
	var v interface{}
	err = tp.Execute(w,v)
	if err != nil {
		fmt.Fprintf(w,"[WebDymap] error %v",err)
		log.Printf("warn: [WebDymap] %s",err)
	}
	log.Printf("[WebDymap INDEX] %s %d %s",r.RemoteAddr,200,r.RequestURI)
}


package web

import (
	"log"
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
//	"github.com/davecgh/go-spew/spew"
)

type WebMain struct {
	Server     *http.Server
	CbStatus   func()(int, string)
	CbStart    func()(int, string)
	TPData     interface{}
}

const Debug = false

func NewWebMain(listen string) *WebMain {
	t := new(WebMain)

	box := boxStatic()

	mux := http.NewServeMux()
	mux.Handle("/static/",http.StripPrefix("/static",http.FileServer(box)))
	mux.HandleFunc("/",t.handleIndex)
	mux.HandleFunc("/api/",t.handleApi)

	t.Server = &http.Server{
		Addr: listen,
		Handler: mux,
	}

	return t
}

func (t *WebMain) Run() {
	log.Printf("info: [WebMain] Start Listening at %s", t.Server.Addr)
	err := t.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (t *WebMain) handleIndex(w http.ResponseWriter, r *http.Request) {
	s, err := boxTemplates("index.tpl.html")
	if err != nil {
		fmt.Fprintf(w,"error")
	}

	tp, err := template.New("T").Parse(s)
	err = tp.Execute(w,t.TPData)

	log.Printf("[WebMain INDEX] %s %d %s",r.RemoteAddr,200,r.RequestURI)
	if err != nil {
		fmt.Fprintf(w,"[WebMain] ERROR %v",err)
		log.Printf("warn: [WebMain] %v",err)
	}
}

func (t *WebMain) handleApi(w http.ResponseWriter, r *http.Request) {

	type apiState struct {
		Code    int    `json:"code"`
		State   int    `json:"state"`
		Message string `json:"message"`
	}
	resp := func(rs apiState) {
		buf,err := json.Marshal(rs)
		if err != nil {
			log.Printf("warn: [Webmain] %s",err)
			w.WriteHeader(200)
			w.Header().Set("Content-Type","application/json; charset=utf-8")
			fmt.Println(w,`{ code: 500, state:0, message:"Internal Server Error" }`)
			log.Printf("%s %s %d",r.RemoteAddr,r.RequestURI,500)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type","application/json; charset=utf-8")
		fmt.Fprintln(w,string(buf))
		if Debug {
			log.Printf("[WebMain API] %s %d %s",r.RemoteAddr,rs.Code,r.RequestURI)
		}
	}

	switch r.URL.Path {
		case "/api/state":
			state,message := t.CbStatus()
			resp(apiState{ Code: 200, State: state, Message: message })
		case "/api/poweron":
			state,message := t.CbStart()
			resp(apiState{ Code: 200, State: state, Message: message })
		default:
			resp(apiState{ Code: 404, State: 0, Message: "File Not Found" })
	}
}


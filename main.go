//go:generate ./build.sh tmpl-dev

// mananno main program
package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mmbros/mananno/httpcache"
	"github.com/mmbros/mananno/jsonrpc"
	"github.com/mmbros/mananno/scraper/arenavision"
	"github.com/mmbros/mananno/templates"
	"github.com/mmbros/mananno/transmission"
)

var (
	httpcacheClient *httpcache.Client
	sch             *arenavision.Schedule
	trans           *transmission.Client
)

func handlerTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := templates.PageTestTransmission.Execute(w, nil); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerRedirect(location string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log.Printf("%s %s", r.Method, r.URL)
		log.Print(r.URL.Query())
		http.Redirect(w, r, location, http.StatusMovedPermanently)
	}
}

func handlerArenavisionSchedule(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	log.Printf("%s %s", r.Method, r.URL)
	log.Print(r.URL.Query())

	err = sch.Get(httpcacheClient)

	events, lastUpdate := sch.Events()
	events = events.FilteredBy(r.URL.Query())
	// filters:
	//   - sport : []string
	//   - competition: []sting
	data := struct {
		HeadTitle  string
		Events     []*arenavision.Event
		LastUpdate time.Time
		Err        error
	}{
		"ArenaVision Schedule",
		events,
		lastUpdate,
		err,
	}

	if err = templates.PageArenavisionSchedule.Execute(w, data); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerArenavisionScheduleRefresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("%s %s", r.Method, r.URL)
	log.Print("*** REFRESH ***")
	httpcacheClient.Clear(sch.SourceURL())
	// redirect to the schedule page, preservng original url query
	newurl := *r.URL
	newurl.Path = "/arenavision/schedule"
	http.Redirect(w, r, newurl.String(), http.StatusMovedPermanently)
}

func handlerArenavisionChannel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("%s %s", r.Method, r.URL)

	channel := arenavision.Channel(ps.ByName("name"))
	link, err := channel.GetLink(httpcacheClient)
	event, live := sch.EventByChannel(channel)

	data := struct {
		HeadTitle string
		Channel   arenavision.Channel
		Stream    template.URL
		Event     *arenavision.Event
		Live      *arenavision.Live
		Err       error
	}{
		channel.FullName(),
		channel,
		template.URL(link),
		event,
		live,
		err,
	}

	if err := templates.PageArenavisionChannel.Execute(w, data); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func jsonrpcSessionGet(req *jsonrpc.Request) (interface{}, error) {

	res, err := trans.SessionGet()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	httpcacheClient = httpcache.NewTTL("/tmp/mananno", 95*time.Minute)
	sch = new(arenavision.Schedule)
	trans = transmission.NewClient("192.168.1.2:9091", "", "")

	//err := sch.Get(httpcacheClient)
	//if err != nil {
	//log.Fatal(err)
	//}

	addr := ":8080"

	// ROUTER
	router := httprouter.New()
	router.GET("/arenavision/schedule", handlerArenavisionSchedule)
	router.GET("/arenavision/schedule/refresh", handlerArenavisionScheduleRefresh)
	router.GET("/arenavision/av:name", handlerArenavisionChannel)

	router.GET("/test", handlerTest)

	router.GET("/", handlerRedirect("/test"))

	// INIT JSON-RPC SERVER
	rpcserver := jsonrpc.NewServer()
	rpcserver.MethodMap["session-get"] = jsonrpcSessionGet
	router.POST("/jsonrpc", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { rpcserver.Handler(w, r) })

	log.Printf("listening to %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Panic(err)
	}
}

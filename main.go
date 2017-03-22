//go:generate ./build.sh tmpl-dev

// mananno main program
package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/mmbros/mananno/httpcache"
	"github.com/mmbros/mananno/jsonrpc"
	"github.com/mmbros/mananno/scraper/arenavision"
	"github.com/mmbros/mananno/scraper/ilcorsaronero"
	"github.com/mmbros/mananno/templates"
	"github.com/mmbros/mananno/transmission"
)

var (
	httpcacheClient *httpcache.Client
	sch             *arenavision.Schedule
	trans           *transmission.Client
)

func handlerCorsaroIndex(w http.ResponseWriter, r *http.Request) {
	var (
		data struct {
			Search   string
			Category ilcorsaronero.Category
			Items    ilcorsaronero.SearchResults
			Err      error
		}
	)

	// check category param
	if i, err := strconv.Atoi(r.FormValue("category")); err != nil {
		data.Category = ilcorsaronero.CatAll
	} else {
		data.Category = ilcorsaronero.Category(i)
	}

	// check search param
	data.Search = r.FormValue("search")

	if len(data.Search) > 0 {
		// do search
		client := ilcorsaronero.Client{}
		data.Items, data.Err = client.Search(data.Search, data.Category)
	}

	if err := templates.PageCorsaroIndex.Execute(w, data); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerTest(w http.ResponseWriter, r *http.Request) {
	if err := templates.PageTestTransmission.Execute(w, nil); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerRedirect(location string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, location, http.StatusMovedPermanently)
	}
}

func handlerArenavisionSchedule(w http.ResponseWriter, r *http.Request) {
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

func handlerArenavisionScheduleRefresh(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL)
	log.Print("*** REFRESH ***")
	httpcacheClient.Clear(sch.SourceURL())
	// redirect to the schedule page, preservng original url query
	newurl := *r.URL
	newurl.Path = "/arenavision/schedule"
	http.Redirect(w, r, newurl.String(), http.StatusMovedPermanently)
}

func handlerArenavisionChannel(w http.ResponseWriter, r *http.Request) {
	ps := FetchParams(r)

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
	return trans.SessionGet()
}

func jsonrpcTorrentAdd(req *jsonrpc.Request) (interface{}, error) {

	type TorrentAddReqParams struct {
		Hash   string `json:"hash"`
		HRef   string `json:"href"`
		Paused bool   `json:"paused"`
	}

	var params TorrentAddReqParams

	err := json.Unmarshal([]byte(*req.Params), &params)
	if err != nil {
		return nil, err
	}
	log.Printf("jsonrpc:torrent-add: hash:%s, paused:%v", params.Hash, params.Paused)

	client := ilcorsaronero.Client{}
	magnet := client.GetMagnet(params.HRef, params.Hash)

	return trans.TorrentAdd(magnet, params.Paused)
}

//func h2hr(fn http.Handler) httprouter.Handle {
//return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//fn(w, r)
//}
//}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		if next != nil {
			next.ServeHTTP(w, r)
		}
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func main() {
	cfg, err := loadConfigFromFile("config.toml")
	if err != nil {
		log.Panic(err)
	}
	log.Print(cfg)
	httpcacheClient = httpcache.NewTTL("/tmp/mananno", 95*time.Minute)
	sch = new(arenavision.Schedule)
	trans = transmission.NewClient(
		cfg.Transmission.Address(),
		cfg.Transmission.Username,
		cfg.Transmission.Password)

	router := httprouter.New()

	commonHandlers := alice.New(loggingHandler)

	// helper functions
	routerGET := func(path string, h http.HandlerFunc) {
		router.GET(path, wrapHandler(commonHandlers.ThenFunc(h)))
	}
	routerPOST := func(path string, h http.HandlerFunc) {
		router.POST(path, wrapHandler(commonHandlers.ThenFunc(h)))
	}
	routerServeFiles := func(path string, root http.FileSystem) {
		if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
			panic("path must end with /*filepath in path '" + path + "'")
		}

		fileServer := http.FileServer(root)

		routerGET(path, func(w http.ResponseWriter, req *http.Request) {
			ps := FetchParams(req)
			req.URL.Path = ps.ByName("filepath")
			fileServer.ServeHTTP(w, req)
		})
	}

	// routes
	routerGET("/arenavision", handlerRedirect("/arenavision/schedule"))
	routerGET("/arenavision/schedule", handlerArenavisionSchedule)
	routerGET("/arenavision/schedule/refresh", handlerArenavisionScheduleRefresh)
	routerGET("/arenavision/av:name", handlerArenavisionChannel)

	routerGET("/ilcorsaronero", handlerCorsaroIndex)
	routerPOST("/ilcorsaronero", handlerCorsaroIndex)

	routerGET("/test", handlerTest)

	routerGET("/", handlerRedirect("/ilcorsaronero"))

	// json-rpc server
	rpcserver := jsonrpc.NewServer()
	rpcserver.MethodMap["session-get"] = jsonrpcSessionGet
	rpcserver.MethodMap["torrent-add"] = jsonrpcTorrentAdd
	routerPOST("/jsonrpc", func(w http.ResponseWriter, r *http.Request) { rpcserver.Handler(w, r) })

	// static files

	//router.ServeFiles("/js/*filepath", http.Dir(cfg.Assets.JS))
	fsJS := &assetfs.AssetFS{
		Asset:     templates.Asset,
		AssetDir:  templates.AssetDir,
		AssetInfo: templates.AssetInfo,
		Prefix:    cfg.Assets.JS,
	}
	routerServeFiles("/js/*filepath", fsJS)

	//router.ServeFiles("/css/*filepath", http.Dir(cfg.Assets.CSS))
	fsCSS := &assetfs.AssetFS{
		Asset:     templates.Asset,
		AssetDir:  templates.AssetDir,
		AssetInfo: templates.AssetInfo,
		Prefix:    cfg.Assets.CSS,
	}
	routerServeFiles("/css/*filepath", fsCSS)

	// start web server
	log.Printf("Starting Mananno web server: listening to %s", cfg.Server.Address())
	if err := http.ListenAndServe(cfg.Server.Address(), router); err != nil {
		log.Panic(err)
	}
}

// ****************************************************************************
// http://www.apriendeau.com/post/middleware-and-httprouter/

func wrap(p string, h func(http.ResponseWriter, *http.Request)) (string, httprouter.Handle) {
	return p, wrapHandler(alice.New(loggingHandler).ThenFunc(h))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

func FetchParams(req *http.Request) httprouter.Params {
	ctx := req.Context()
	return ctx.Value("params").(httprouter.Params)
}

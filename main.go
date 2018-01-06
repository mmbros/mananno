//go:generate ./build.sh tmpl-dev

// mananno main program
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/mmbros/mananno/httpcache"
	"github.com/mmbros/mananno/jsonrpc"
	"github.com/mmbros/mananno/scraper/acestreamid"
	"github.com/mmbros/mananno/scraper/arenavision"
	"github.com/mmbros/mananno/scraper/ilcorsaronero"
	"github.com/mmbros/mananno/scraper/tntvillage"
	"github.com/mmbros/mananno/templates"
	"github.com/mmbros/mananno/transmission"
)

var (
	cfg             *configuration
	httpcacheClient *httpcache.Client
	av              *arenavision.Scraper
	scprAcestreamid *acestreamid.Scraper
	trans           *transmission.Client
)

func handlerHomepage(w http.ResponseWriter, r *http.Request) {
	if err := templates.PageHomepage.Execute(w, nil); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerTntVillageIndex(w http.ResponseWriter, r *http.Request) {
	var (
		data struct {
			Search   string
			Category tntvillage.Category
			Items    tntvillage.SearchResults
			Err      error
		}
	)

	// check category param
	if i, err := strconv.Atoi(r.FormValue("category")); err != nil {
		data.Category = tntvillage.CatAll
	} else {
		data.Category = tntvillage.Category(i)
	}

	// check search param
	data.Search = r.FormValue("search")

	if len(data.Search) > 0 {
		// do search
		client := tntvillage.Client{}
		data.Items, data.Err = client.Search(data.Search, data.Category)
	}

	if err := templates.PageTntVillageIndex.Execute(w, data); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerTransmission(w http.ResponseWriter, r *http.Request) {
	ps := FetchParams(r)
	filepath := ps.ByName("filepath")
	log.Print("filepath: ", filepath)
	log.Print("Transmission.Web: ", cfg.Transmission.Web())
	u, _ := url.Parse(cfg.Transmission.URL())
	u.Path = path.Join(u.Path, filepath)
	s := u.String()
	if strings.HasSuffix(filepath, "/") {
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
	}
	log.Print("URL: ", s)

	resp, err := trans.Get(s)
	if err != nil {
		log.Printf("!!! Get error !!!\n")
		log.Fatal(err)
	}

	if strings.HasSuffix(filepath, ".css") {
		log.Print("CSS:", filepath)
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}

	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "\n\n%s\n\n", text)
	if err != nil {
		log.Fatal(err)
	}
}
func handlerTransmissionRevProxy(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-Go-Proxy", "MMbros")
		p.ServeHTTP(w, r)
	}
}

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

	err = av.RefreshGuide()

	lastUpdate := time.Now()
	events := av.Events
	events = events.FilteredBy(r.URL.Query())
	// filters:
	//   - sport : []string
	//   - competition: []sting
	data := struct {
		HeadTitle  string
		Events     arenavision.Events
		LastUpdate time.Time
		Err        error
	}{
		"ArenaVision Events Guide",
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
	httpcacheClient.Clear(av.GuideURL())
	// redirect to the schedule page, preserving original url query
	newurl := *r.URL
	newurl.Path = "/arenavision/schedule"
	http.Redirect(w, r, newurl.String(), http.StatusMovedPermanently)
}

func handlerArenavisionChannel(w http.ResponseWriter, r *http.Request) {
	ps := FetchParams(r)

	channel := av.Channels[ps.ByName("name")]
	link, err := channel.GetLink(httpcacheClient)
	event, live := av.EventByChannel(channel)

	data := struct {
		HeadTitle string
		Channel   *arenavision.Channel
		Stream    template.URL
		Event     *arenavision.Event
		Live      *arenavision.Live
		Err       error
	}{
		channel.Name,
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
func handlerAcestreamidChannels(w http.ResponseWriter, r *http.Request) {
	var err error
	log.Printf("%s %s", r.Method, r.URL)
	log.Print(r.URL.Query())

	err = scprAcestreamid.Refresh()
	log.Printf("Acestreamid channels: %d\n", len(scprAcestreamid.Channels))

	lastUpdate := time.Now()

	data := struct {
		HeadTitle  string
		Channels   acestreamid.Channels
		LastUpdate time.Time
		Err        error
	}{
		"Acestreamid CANALI",
		scprAcestreamid.Channels,
		lastUpdate,
		err,
	}

	if err = templates.PageAcestreamidChannels.Execute(w, data); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}

func handlerAcestreamidChannel(w http.ResponseWriter, r *http.Request) {
	var err error
	log.Printf("%s %s", r.Method, r.URL)
	log.Print(r.URL.Query())

	ps := FetchParams(r)

	channel := scprAcestreamid.ChannelByID(ps.ByName("name"))
	log.Printf("CHANNEL %s: %v\n", ps.ByName("name"), channel)
	if channel == nil {
		http.NotFound(w, r)
		return
	}
	channel.Refresh(scprAcestreamid.Client)

	lastUpdate := time.Now()

	data := struct {
		HeadTitle  string
		Name       string
		Streams    acestreamid.Streams
		LastUpdate time.Time
		Err        error
	}{
		channel.ID(),
		channel.Name,
		channel.Streams,
		lastUpdate,
		err,
	}

	if err = templates.PageAcestreamidChannel.Execute(w, data); err != nil {
		log.Printf("Template error: %q\n", err)
	}
}
func jsonrpcSessionGet(req *jsonrpc.Request) (interface{}, error) {
	return trans.SessionGet()
}

func jsonrpcMagnetAdd(req *jsonrpc.Request) (interface{}, error) {

	type MagnetAddReqParams struct {
		Magnet string `json:"magnet"`
		Paused bool   `json:"paused"`
	}

	var params MagnetAddReqParams

	err := json.Unmarshal([]byte(*req.Params), &params)
	if err != nil {
		return nil, err
	}
	log.Printf("jsonrpc:magnet-add: magnet:%s, paused:%v", params.Magnet, params.Paused)

	return trans.TorrentAdd(params.Magnet, params.Paused)
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
	var err error
	cfg, err = loadConfigFromFile("config.toml")
	if err != nil {
		log.Panic(err)
	}
	log.Print(cfg)
	httpcacheClient = httpcache.NewTTL("/tmp/mananno", 95*time.Minute)
	av = &arenavision.Scraper{Client: httpcacheClient}
	scprAcestreamid = &acestreamid.Scraper{Client: httpcacheClient}

	trans = transmission.NewClient(
		cfg.Transmission.Address(),
		cfg.Transmission.Username,
		cfg.Transmission.Password)

	router := httprouter.New()

	commonHandlers := alice.New(loggingHandler)

	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:     templates.Asset,
		AssetDir:  templates.AssetDir,
		AssetInfo: templates.AssetInfo,
	})

	// helper functions
	routerGET := func(path string, h http.HandlerFunc) {
		router.GET(path, wrapHandler(commonHandlers.ThenFunc(h)))
	}
	routerPOST := func(path string, h http.HandlerFunc) {
		router.POST(path, wrapHandler(commonHandlers.ThenFunc(h)))
	}

	// routes
	routerGET("/", handlerHomepage)

	routerGET("/arenavision", handlerRedirect("/arenavision/schedule"))
	routerGET("/arenavision/schedule", handlerArenavisionSchedule)
	routerGET("/arenavision/schedule/refresh", handlerArenavisionScheduleRefresh)
	routerGET("/arenavision/channels/:name", handlerArenavisionChannel)

	routerGET("/tntvillage", handlerTntVillageIndex)
	routerPOST("/tntvillage", handlerTntVillageIndex)

	routerGET("/ilcorsaronero", handlerCorsaroIndex)
	routerPOST("/ilcorsaronero", handlerCorsaroIndex)

	routerGET("/acestreamid", handlerRedirect("/acestreamid/channels"))
	routerGET("/acestreamid/channels", handlerAcestreamidChannels)
	routerGET("/acestreamid/channels/:name", handlerAcestreamidChannel)

	routerGET("/test", handlerTest)

	// json-rpc server
	rpcserver := jsonrpc.NewServer()
	rpcserver.MethodMap["session-get"] = jsonrpcSessionGet
	rpcserver.MethodMap["torrent-add"] = jsonrpcTorrentAdd
	rpcserver.MethodMap["magnet-add"] = jsonrpcMagnetAdd
	//routerPOST("/jsonrpc", func(w http.ResponseWriter, r *http.Request) { rpcserver.Handler(w, r) })
	routerPOST("/jsonrpc", rpcserver.Handler)

	// static files
	routerGET("/css/*filepath", fileServer.ServeHTTP)
	routerGET("/img/*filepath", fileServer.ServeHTTP)
	routerGET("/js/*filepath", fileServer.ServeHTTP)

	// transmission reverse proxy
	//log.Print(cfg.Transmission.Web())
	//remote, err := url.Parse(cfg.Transmission.Web())
	//if err != nil {
	//panic(err)
	//}
	//proxy := httputil.NewSingleHostReverseProxy(remote)
	//routerGET("/transmission/*filepath", handlerTransmissionRevProxy(proxy))

	routerGET("/transmission/*filepath", handlerTransmission)

	// start web server
	log.Printf("Starting Mananno web server: listening to %s", cfg.Server.Address())
	if err := http.ListenAndServe(cfg.Server.Address(), router); err != nil {
		log.Panic(err)
	}
}
func main3() {
	cfg, err := loadConfigFromFile("config.toml")
	if err != nil {
		log.Panic(err)
	}

	trans = transmission.NewClient(
		cfg.Transmission.Address(),
		cfg.Transmission.Username,
		cfg.Transmission.Password,
	)

	resp, err := trans.Get(cfg.Transmission.Web())
	if err != nil {
		log.Printf("!!! Get error !!!\n")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("\n\n%s\n\n", text)
	if err != nil {
		log.Fatal(err)
	}

}

// ****************************************************************************
// http://www.apriendeau.com/post/middleware-and-httprouter/
// https://blog.golang.org/context#TOC_3.2.

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// paramsKey is the context key for the httprouter.Params. Its value of zero is
// arbitrary. If this package defined other context keys, they would have
// different integer values.
const paramsKey contextKey = 0

func wrap(p string, h func(http.ResponseWriter, *http.Request)) (string, httprouter.Handle) {
	return p, wrapHandler(alice.New(loggingHandler).ThenFunc(h))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), paramsKey, ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

// FetchParams returns the httprouter.Params of the given http.Request.
func FetchParams(req *http.Request) httprouter.Params {
	ctx := req.Context()
	return ctx.Value(paramsKey).(httprouter.Params)
}

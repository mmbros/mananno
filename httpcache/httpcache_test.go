package httpcache

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"testing"
	"time"
)

var cacheFolder = "/tmp/httpcache_test"
var logbuf bytes.Buffer

func init() {
	log.SetOutput(&logbuf)
}

func myTestDumpResponse(t *testing.T) {

	filename := "/tmp/test_dump_response.html"
	url := `https://golang.org`

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Read Dumped Response

	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	buf2 := bufio.NewReader(f)
	resp2, err := http.ReadResponse(buf2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.Status != "200 OK" {
		t.Errorf("Expecting status 200 OK, found %s\n", resp2.Status)
	}
}
func TestArenavisionSchedule(t *testing.T) {

	url := "http://arenavision.in/schedule"
	ttl := 10 * time.Minute
	client := NewTTL(cacheFolder, ttl)

	// create a new http.Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// set beget cookie
	cookie := http.Cookie{
		Name:    "beget",
		Value:   "begetok",
		Path:    "/",
		Expires: time.Now().Add(19360000000),
	}
	req.AddCookie(&cookie)

	// get the response
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestGet(t *testing.T) {

	checklog := func(msg string) {
		if !strings.Contains(logbuf.String(), msg) {
			t.Log(logbuf.String())
			t.Fatalf("Expecting %q", msg)
		}
		logbuf.Reset()
	}

	url := "http://www.google.com"
	ttl := 100 * time.Millisecond
	client := NewTTL(cacheFolder, ttl)

	// clean cached file, if exists
	os.Remove(client.LocalPath(url))

	// GET 1: get from net
	resp, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	checklog("cached file not exists")

	// GET 2: get from cache
	resp2, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp2.Body.Close()
	checklog("using cached file")

	// GET 3: cache expired
	time.Sleep(ttl)
	resp3, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp3.Body.Close()
	checklog("cached file expired")
}

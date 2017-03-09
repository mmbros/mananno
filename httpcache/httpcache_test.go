package httpcache

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var logbuf bytes.Buffer

func init() {
	log.SetOutput(&logbuf)
}

//func TestDumpResponse(t *testing.T) {

//filename := "/tmp/test_dump_response.html"
//url := `https://golang.org`

//resp, err := http.Get(url)
//if err != nil {
//t.Fatal(err)
//}
//defer resp.Body.Close()

//buf, err := httputil.DumpResponse(resp, true)
//if err != nil {
//t.Fatal(err)
//}

//err = ioutil.WriteFile(filename, buf, 0644)
//if err != nil {
//t.Fatal(err)
//}

//// Read Dumped Response

//f, err := os.Open(filename)
//if err != nil {
//t.Fatal(err)
//}
//defer f.Close()

//buf2 := bufio.NewReader(f)
//resp2, err := http.ReadResponse(buf2, nil)
//if err != nil {
//t.Fatal(err)
//}
//defer resp2.Body.Close()

//if resp2.Status != "200 OK" {
//t.Errorf("Expecting status 200 OK, found %s\n", resp2.Status)
//}
//}

func TestGet(t *testing.T) {

	checklog := func(msg string) {
		if !strings.Contains(logbuf.String(), msg) {
			t.Log(logbuf.String())
			t.Fatalf("Expecting %q", msg)
		}
		logbuf.Reset()
	}

	folder := "/tmp/httpcache_test"
	url := "http://www.google.com"
	ttl := 100 * time.Millisecond
	client := NewTTL(folder, ttl)

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

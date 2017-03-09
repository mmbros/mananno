package httpcache

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
	"time"
)

func TestDumpResponse(t *testing.T) {

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
func TestGet(t *testing.T) {
	folder := "/tmp/httpcache_test"
	url := "https://golang.org"
	client := NewTTL(folder, 2*time.Second)

	resp, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	resp2, err := client.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp2.Body.Close()
}

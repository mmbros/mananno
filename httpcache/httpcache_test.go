package httpcache

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
)

func TestDumpResponse(t *testing.T) {

	filename := "/tmp/test_dump_response.html"
	url := `http://arenavision.in/av1`
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}

	err = ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}

	// Read Dumped Response

	f, err := os.Open(filename)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}
	defer f.Close()

	buf2 := bufio.NewReader(f)
	resp2, err := http.ReadResponse(buf2, nil)
	if err != nil {
		t.Errorf("Error: %v\n", err)
		return
	}
	defer resp2.Body.Close()

	if resp2.Status != "200 OK" {
		t.Errorf("Expecting status 200 OK, found %s\n", resp2.Status)
	}

}

package httpcache

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"
)

// MaxAger function returns the TTL duration for the given url.
type MaxAger func(url string) time.Duration

// Client is an HTTP client that cached responses.
type Client struct {
	http.Client

	// CacheFolder is the folder where the responses are saved in.
	// If the folder doesn't exist, it is created only when needed.
	CacheFolder string

	// MaxAge function returns the TTL duration for each url.
	MaxAge MaxAger
}

func defaultMaxAgeFactory(ttl time.Duration) MaxAger {
	return func(string) time.Duration {
		return ttl
	}
}

// NewTTL creates a new httpcache.Client that caches all the responses
// for the same constant ttl duration.
func NewTTL(folder string, ttl time.Duration) *Client {
	client := &Client{
		Client: http.Client{
			// Don’t use Go’s default HTTP client (in production)
			// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.q5iexu8v7
			Timeout: 60 * time.Second,
		},
		MaxAge:      defaultMaxAgeFactory(ttl),
		CacheFolder: folder,
	}
	return client
}

// Hash returns the hash used for the given url.
func (client *Client) Hash(url string) string {
	// get the base64 hash of the url
	return base64.StdEncoding.EncodeToString([]byte(url))
}

// LocalPath returns the path of the cached file
func (client *Client) LocalPath(url string) string {
	filename := client.Hash(url) + ".html"
	path := filepath.Join(client.CacheFolder, filename)
	return path
}

func readCachedResponse(url, filename string) (*http.Response, error) {
	// Read Dumped Response
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// return http.ReadResponse(buf, req)
	resp, err := http.ReadResponse(buf, req)
	if err != nil {
		return nil, err
	}

	// save response body
	// see: http://stackoverflow.com/questions/33963467/parse-http-requests-and-responses-from-text-file-in-go?rq=1
	b := new(bytes.Buffer)
	io.Copy(b, resp.Body)
	resp.Body.Close()
	resp.Body = ioutil.NopCloser(b)

	return resp, nil
}

func writeBuffer(filename string, data []byte) error {
	log.Printf("client.WriteBuffer(filename=%q)\n", filename)
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		// in caso di errore prova a creare la cartella e ci riprova
		if os.IsNotExist(err) {
			log.Printf("  WB: creating folder: %s\n", filepath.Dir(filename))
			err2 := os.MkdirAll(filepath.Dir(filename), 0777)
			if err2 == nil {
				err = ioutil.WriteFile(filename, data, 0644)
			} else {
				log.Printf("  ERR creating folder: %q\n", err)
			}
		}
	}
	if err != nil {
		log.Printf("  ERR client.WriteBuffer: %q\n", err)
	}
	return err
}

// Do sends an HTTP request and returns an HTTP response, following policy
// (such as redirects, cookies, auth) as configured on the client.
func (client *Client) Do(req *http.Request) (*http.Response, error) {
	// aux function to return error
	ko := func(err error) (*http.Response, error) {
		log.Printf("  ERR: %v\n", err)
		return nil, err
	}

	// get the url
	url := req.URL.String()
	log.Printf("httpcache.Do(req.url=\"%s\")\n", url)
	// get the filename
	filename := client.LocalPath(url)
	// check if filename exists
	fileinfo, err := os.Stat(filename)
	if err != nil {
		// check os.Stat error
		if !os.IsNotExist(err) {
			// error is not file not exists
			return ko(err)
		}
		// error is file not exists
		log.Printf("  cached file not exists: %s\n", filename)
	} else {
		// file found in cache: check age
		ttl := client.MaxAge(url)
		elapsed := time.Since(fileinfo.ModTime())
		log.Printf("  ttl=%q, elapsed=%q\n", ttl, elapsed)

		if elapsed > ttl {
			log.Printf("  cached file expired: %v\n", fileinfo.ModTime())
			err := os.Remove(filename)
			if err != nil {
				return ko(err)
			}
		} else {
			log.Printf("  using cached file: %s\n", filename)
			return readCachedResponse(url, filename)
		}
	}

	// cache not found or expired

	// get the data from the net using http.Client
	log.Printf("  get url: %s\n", url)
	resp, err := client.Client.Do(req)
	if err != nil {
		return ko(err)
	}

	// dump the resp to buffer
	log.Println("  dumping response")
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return ko(err)
	}

	//log.Printf("  creating cached file: %s\n", filename)
	writeBuffer(filename, buf)

	return resp, nil

}

// Get issues a GET to the specified URL.
func (client *Client) Get(url string) (*http.Response, error) {
	log.Printf("httpcache.Get(url=\"%s\")\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

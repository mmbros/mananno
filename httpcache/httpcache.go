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
	"sync"
	"time"
)

// MaxAger function returns the TTL duration for the given url.
type MaxAger func(url string) time.Duration

type Cache struct {
	folder string
	mx     sync.Mutex
	maxAge MaxAger
	client *http.Client
}

func defaultMaxAgeFactory(ttl time.Duration) MaxAger {
	return func(string) time.Duration {
		return ttl
	}
}

// New creates a new cache client.
// The responses are saved in the given folder. If the folder doesn't exist,
// it is created only when a response must be saved.
// If fnMaxAge is nil, it is used a constant TTL duration of 60 seconds for
// every response.
// If client is nil, a default http.Client with 10 seconds timeout is used.
func New(client *http.Client, folder string, fnMaxAge MaxAger) *Cache {
	if client == nil {
		// Don’t use Go’s default HTTP client (in production)
		// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.q5iexu8v7
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	if fnMaxAge == nil {
		fnMaxAge = defaultMaxAgeFactory(60 * time.Second)
	}
	cache := &Cache{
		folder: folder,
		mx:     sync.Mutex{},
		maxAge: fnMaxAge,
		client: client,
	}
	return cache
}

// NewTTL creates a cache client with constant ttl duration.
// See New for other parameters.
func NewTTL(client *http.Client, folder string, ttl time.Duration) *Cache {
	fn := defaultMaxAgeFactory(ttl)
	return New(client, folder, fn)
}

// Folder returns the cache folder.
func (cache *Cache) Folder() string {
	return cache.folder
}

// Hash returns the hash used for the given url.
func (cache *Cache) Hash(url string) string {
	// get the base64 hash of the url
	return base64.StdEncoding.EncodeToString([]byte(url))
}

// LocalPath returns the path of the cached file
func (cache *Cache) LocalPath(url string) string {
	filename := cache.Hash(url) + ".html"
	path := filepath.Join(cache.folder, filename)
	return path
}

func (cache *Cache) Get(url string) (*http.Response, error) {
	// aux function to return error
	ko := func(err error) (*http.Response, error) {
		log.Printf("  ERR: %v\n", err)
		return nil, err
	}
	log.Printf("cache.Get(url=\"%s\")\n", url)
	// get the filename
	filename := cache.LocalPath(url)
	// check if filename exists
	fileinfo, err := os.Stat(filename)
	if err != nil {
		// file not found in cache
		if !os.IsNotExist(err) {
			return ko(err)
		}
		log.Printf("  cached file not exists: %s\n", filename)
	} else {
		// file found un cache: check age
		ttl := cache.maxAge(url)
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

	// get the data from the net
	log.Printf("  get url: %s\n", url)
	resp, err := cache.client.Get(url)
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
	//return http.ReadResponse(buf, req)
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
	log.Printf("cache.WriteBuffer(filename=%q)\n", filename)
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
		log.Printf("  ERR cache.WriteBuffer: %q\n", err)
	}
	return err
}

func (client *Cache) Do(req *http.Request) (*http.Response, error) {
	// aux function to return error
	ko := func(err error) (*http.Response, error) {
		log.Printf("  ERR: %v\n", err)
		return nil, err
	}

	// get the url
	url := req.URL.String()
	log.Printf("cache.Do(req.url=\"%s\")\n", url)
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
		ttl := client.maxAge(url)
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

	// get the data from the net
	log.Printf("  get url: %s\n", url)
	resp, err := client.client.Do(req)
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

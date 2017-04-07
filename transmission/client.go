package transmission

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Client represents a Transmission client
type Client struct {
	// Address for Json request
	Address  string
	Username string
	Password string
	session  session
	client   *http.Client
}

const headerTransmissionSessionID = "X-Transmission-Session-Id"

type session struct {
	sync.RWMutex
	id string
}

func (s *session) Set(id string) {
	s.Lock()
	log.Printf("session.Set(%s)\n", id)
	s.id = id
	s.Unlock()
}

func (s *session) Get() string {
	s.RLock()
	defer s.RUnlock()
	log.Printf("session.Get() -> %s\n", s.id)
	return s.id
}

func normalizeAddress(addr string) string {
	if u, err := url.Parse(addr); err == nil {
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		// redo parse to properly handle addresses with format "ip:port"
		u, _ := url.Parse(u.String())
		if u.Path == "" {
			u.Path = "/transmission/rpc"
		}
		log.Printf("transmission: updating Address from %q to %q", addr, u.String())
		addr = u.String()
	}
	return addr
}

// NewClient initialize a new Transmission client
func NewClient(addr, usr, pwd string) *Client {
	client := &http.Client{
		// Don’t use Go’s default HTTP client (in production)
		// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.q5iexu8v7
		Timeout: 60 * time.Second,
	}

	return &Client{
		Address:  normalizeAddress(addr),
		Username: usr,
		Password: pwd,
		client:   client,
	}
}

// Ping check the connection
func (c *Client) Ping() bool {
	_, err := c.SessionGet()
	return err == nil
}

// exec function
func (c *Client) exec(method string, args interface{}, reply interface{}) error {

	const maxIter = 3

	// json encoded request arguments
	buf, err := encodeClientRequest(method, &args)
	if err != nil {
		return err
	}

	// get saved X-Transmission-Session-Id
	sessionID := c.session.Get()

	for iter := 0; iter < maxIter; iter++ {
		log.Printf("iter %d: %s\n", iter, string(buf))

		body := bytes.NewBuffer(buf)

		// prepare the request
		req, err := http.NewRequest(http.MethodPost, c.Address, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		if sessionID != "" {
			req.Header.Set(headerTransmissionSessionID, sessionID)
		}
		req.SetBasicAuth(c.Username, c.Password)

		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}

		log.Printf("iter %d: %s\n", iter, resp.Status)
		switch resp.StatusCode {
		case http.StatusConflict:
			// handle "409 Conflict"

			// read and save new X-Transmission-Session-Id
			sessionID = resp.Header[headerTransmissionSessionID][0]
			c.session.Set(sessionID)
			// repeat the request
			continue
		case http.StatusOK:
			// handle "200 OK"
			defer resp.Body.Close()

			if reply != nil {
				err = decodeClientResponse(resp.Body, &reply)
				if err != nil {
					return err
				}
			}
			return nil

		default:
			log.Printf("status not expected: %s\n", resp.Status)
			log.Print(resp.Header)
			return fmt.Errorf("Transmission: %s", resp.Status)
		}
	}

	return errors.New("exec: too many iter")
}

// TorrentAddResponse represents the information returned
// by the TorrentAdd method.
type TorrentAddResponse struct {
	TorrentAdded struct {
		ID         int    `json:"id"`
		HashString string `json:"hashString"`
		Name       string `json:"name"`
	} `json:"torrent-added"`
}

// TorrentAdd method
func (c *Client) TorrentAdd(URL string, paused bool) (*TorrentAddResponse, error) {
	const method = "torrent-add"
	// torrent-add request arguments
	args := struct {
		Filename string `json:"filename"`
		Paused   bool   `json:"paused"`
	}{
		URL,
		paused,
	}

	var reply TorrentAddResponse

	err := c.exec(method, &args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// TorrentRemove method
func (c *Client) TorrentRemove(ids string, deleteLocalData bool) error {
	const method = "torrent-remove"

	// request arguments
	args := struct {
		Ids             string `json:"ids"`
		DeleteLocalData bool   `json:"delete-local-data"`
	}{
		ids,
		deleteLocalData,
	}

	//err := c.exec(method, &args, nil)
	//if err != nil {
	//return err
	//}
	//return nil
	return c.exec(method, &args, nil)
}

// SessionGetResponse represents the information returned
// by the SessionGet method.
type SessionGetResponse struct {
	// location of transmission's configuration directory
	ConfigDir string `json:"config-dir"`
	// default path to download torrents
	DownloadDir string `json:"download-dir"`
	//  max number of torrents to download at once (see download-queue-enabled)
	DownloadQueueSize int `json:"download-queue-size"`
	// if true, limit how many torrents can be downloaded at once
	DownloadQueueEnabled bool `json:"download-queue-enabled"`
	// true means allow dht in public torrents
	DhtEnabled bool `json:"dht-enabled"`
	// "required", "preferred", "tolerated"
	Encryption string
	// torrents we're seeding will be stopped if they're idle for this long
	IdleSeedingLimit int `json:"idle-seeding-limit"`
	// true if the seeding inactivity limit is honored by default
	IdleSeedingLimitEnabled bool `json:"idle-seeding-limit-enabled"`
	// path for incomplete torrents, when enabled
	InclompleteDir string `json:"incomplete-dir"`
	// true means keep torrents in incomplete-dir until done
	InclompleteDirEnabled bool `json:"incomplete-dir-enabled"`
	// the current RPC API version
	RPCVersion int `json:"rpc-version"`
	// the minimum RPC API version supported
	RPCVersionMinimum int `json:"rpc-version-minimum"`
	// long version string "$version ($revision)"
	Version string
}

// SessionGet method
func (c *Client) SessionGet() (*SessionGetResponse, error) {
	const method = "session-get"

	var reply SessionGetResponse

	err := c.exec(method, nil, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// clientRequest represents a JSON-RPC request sent by a client.
type clientRequest struct {
	// A required "method" string telling the name of the method to invoke
	Method string `json:"method"`

	// An optional "arguments" object of key/value pairs
	Args interface{} `json:"arguments"`

	// An optional "tag" number used by clients to track responses.
	// If provided by a request, the response MUST include the same tag.
	Tag uint64 `json:"tag"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	// A required "result" string whose value MUST be "success" on success,
	// or an error string on failure.
	Result string `json:"result"`

	// An optional "arguments" object of key/value pairs
	Args *json.RawMessage `json:"arguments"`

	// An optional "tag" number as described in 2.1.
	Tag uint64 `json:"tag"`
}

// encodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		Method: method,
		Args:   args,
		Tag:    uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

// decodeClientResponse decodes the response body of a client request into
// the interface reply.
func decodeClientResponse(r io.Reader, reply interface{}) error {
	var c clientResponse
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}
	if c.Result != "success" {
		return errors.New(c.Result)
	}

	return json.Unmarshal(*c.Args, reply)
}

// Get function
func (c *Client) Get(url string) (*http.Response, error) {
	log.Printf("transmission.Get(%q)", url)

	const maxIter = 3

	// get saved X-Transmission-Session-Id
	sessionID := c.session.Get()

	for iter := 0; iter < maxIter; iter++ {

		// prepare the request
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		if sessionID != "" {
			req.Header.Set(headerTransmissionSessionID, sessionID)
		}
		req.SetBasicAuth(c.Username, c.Password)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}

		log.Printf("iter %d: %s\n", iter, resp.Status)
		switch resp.StatusCode {
		case http.StatusConflict:
			// handle "409 Conflict"

			// read and save new X-Transmission-Session-Id
			sessionID = resp.Header[headerTransmissionSessionID][0]
			c.session.Set(sessionID)
			// repeat the request
			continue
		case http.StatusOK:
			// handle "200 OK"
			return resp, nil

		default:
			log.Printf("status not expected: %s\n", resp.Status)
			log.Print(resp.Header)
			return nil, fmt.Errorf("Transmission: %s", resp.Status)
		}
	}

	return nil, errors.New("exec: too many iter")
}

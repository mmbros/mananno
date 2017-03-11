package transmission

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// TorrentAddResponse represents the information returned
// by the TorrentAdd method.
type TorrentAddResponse struct {
	TorrentAdded struct {
		ID         int    `json:"id"`
		HashString string `json:"hashString"`
		Name       string `json:"name"`
	} `json:"torrent-added"`
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

// Client represents a Transmission client
type Client struct {
	Address  string
	Username string
	Password string
	session  session
}

// NewClient initialize a new Transmission client
func NewClient(addr, usr, pwd string) *Client {
	return &Client{
		Address:  addr,
		Username: usr,
		Password: pwd}
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

	client := &http.Client{
		// Don’t use Go’s default HTTP client (in production)
		// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.q5iexu8v7
		Timeout: 60 * time.Second,
	}

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

		//log.Println("BEFORE client.Do")
		resp, err := client.Do(req)
		//log.Println("AFTER client.Do")
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
		}
	}

	return errors.New("exec: too many iter")
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

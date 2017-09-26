package arenavision

import (
	"strings"

	"github.com/mmbros/www/scraper"
)

// Channel represents an Arenavision channel.
type Channel struct {
	Name string
	URL  string
}

// Channels is a map of ArenaVision channels.
type Channels map[string]*Channel

func (chs Channels) Get(s string) *Channel {
	s = strings.TrimSpace(s)
	ch := chs[s]
	if ch == nil {
		ch = chs["0"+s]
	}
	return ch
}

func (ch *Channel) SourceURL() string {
	return absURL(ch.URL)
}

// GetLink get the channel source page and returns the first stream link found.
// Limits: currently it handles only acestream links.
func (ch *Channel) GetLink(client scraper.URLGetter) (string, error) {

	resp, err := getURL(client, ch.SourceURL())
	if err != nil {
		return "", err
	}
	return scraper.GetFirstAcestreamLink(resp)
}

// ID returns the identifier of an Arenavision channel.
func (ch *Channel) ID() string {
	u := ch.URL
	if (len(u) == 0) || (u[0] != '/') {
		return u
	}
	return u[1:]
}

func (ch *Channel) String() string {
	return ch.ID()
}

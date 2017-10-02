package acestreamid

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/mananno/scraper"
)

// Scraper is the Acestreamid scraper.
type Scraper struct {
	Channels   Channels
	Client     scraper.URLGetter
	lastUpdate time.Time
}

// Channel represents an Acestreamid.com channel information.
type Channel struct {
	Name  string
	Href  string
	Count string
}

// Channels is a collenction of Acestreamid.com Channel.
type Channels map[string]*Channel

// Stream represents an Acestreamid.com stream information.
type Stream struct {
	Title string
	ID    string
	Site  string
	Time  string
}

// Streams is an array of Acestreamid.com Stream.
type Streams []*Stream

// ID returns the identifier of the channel
func (ch *Channel) ID() string {
	return strings.TrimPrefix(ch.Href, "/channel/")
}
func parseStreams(doc *goquery.Document) (Streams, error) {
	streams := Streams{}
	doc.Find("li.collection-item").Each(func(i int, s *goquery.Selection) {
		strm := &Stream{
			Title: s.Find(".col_title").Text(),
			ID:    s.Find(".col_id").Text(),
			Time:  strings.TrimSpace(s.Find(".col_time span").Text()),
			Site:  s.Find(".col_time a").AttrOr("href", ""),
		}
		streams = append(streams, strm)
	})
	return streams, nil
}

func parseChannels(doc *goquery.Document) (Channels, error) {
	channels := Channels{}
	doc.Find("li.collection-item").Each(func(i int, s *goquery.Selection) {
		link := s.Find(".link")
		ch := &Channel{
			Name:  link.Text(),
			Count: strings.TrimSpace(s.Find(".content").Text()),
		}
		if ch.Name != "" {
			ch.Href, _ = link.Attr("href")
			channels[ch.ID()] = ch
		}
	})
	return channels, nil
}

func (scpr *Scraper) client() scraper.URLGetter {
	if scpr.Client != nil {
		return scpr.Client
	}
	return scraper.DefaultURLGetter()
}
func getURL(client scraper.URLGetter, URL string) (*http.Response, error) {
	return client.Get(URL)
}

// Refresh updates the scraper informations.
func (scpr *Scraper) Refresh() error {
	u := "https://acestreamid.com/stat/channels"
	resp, err := getURL(scpr.client(), u)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// create a goquery document from the response
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	// get the urls of the channels
	scpr.Channels, _ = parseChannels(doc)

	return nil
}

// ChannelByID returns the channel with given id.
// It returns nil if the channel is not found.
func (scpr *Scraper) ChannelByID(id string) *Channel {
	if len(scpr.Channels) == 0 {
		scpr.Refresh()
	}
	return scpr.Channels[id]
}

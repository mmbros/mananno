package acestreaid

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Channel struct {
	Name  string
	Href  string
	Count string
}

// Channels is an array of Acestreamid.com Channel.
type Channels []*Channel

// Stream represents an Acestreamid.com stream information.
type Stream struct {
	Title string
	ID    string
	Site  string
	Time  string
}

// Streams is an array of Acestreamid.com Stream.
type Streams []*Stream

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
		ch.Href, _ = link.Attr("href")
		channels = append(channels, ch)
	})
	return channels, nil
}

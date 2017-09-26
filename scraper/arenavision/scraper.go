package arenavision

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/www/scraper"
)

const (
	baseURL         = "http://arenavision.in/"
	defaultGuideURL = "/iguide"
)

// Arenavision scraper.
type Scraper struct {
	Channels   Channels
	Events     Events
	Client     scraper.URLGetter
	lastUpdate time.Time
	guideURL   string // use defaultGuideURL if empty
}

func (av *Scraper) client() scraper.URLGetter {
	if av.Client != nil {
		return av.Client
	}
	return scraper.DefaultURLGetter()
}

func (av *Scraper) GuideURL() string {
	u := av.guideURL
	if u == "" {
		u = defaultGuideURL
	}
	return absURL(u)
}

func (av *Scraper) refreshGuideURL() error {
	// get arenavision homepage
	resp, err := getURL(av.client(), baseURL)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return errors.New("Arenavision homepage not found!")
	}

	defer resp.Body.Close()

	// create a goquery document from the response
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	// get the url of events guide
	href := parseGuideURL(doc)
	if href == "" {
		return errors.New("Events guide URL not found")
	}

	// update guideURL
	av.guideURL = href

	// get the urls of the channels
	// av.Channels = parseChannels(doc)

	return nil
}

func (av *Scraper) RefreshGuide() error {
	firstIter := true
start:
	// get arenavision guide url
	u := av.GuideURL()
	resp, err := getURL(av.client(), u)
	if (err != nil) || (resp.StatusCode != http.StatusOK) {
		if av.guideURL == "" {
			// default guide url was used.
			// get the actual guide url from homepage.
			if err := av.refreshGuideURL(); err != nil {
				return err
			}
			// retry to get the actual guide url
			if firstIter {
				firstIter = false
				goto start
			}
		}
		return errors.New("Can't get Events Guide page")
	}

	defer resp.Body.Close()

	// create a goquery document from the response
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	// get the urls of the channels
	av.Channels = parseChannels(doc)

	// get the scheduled events
	av.Events = parseEvents(doc)
	for _, e := range av.Events {
		e.Lives = stringToLives(e.live, av.Channels)

	}

	return nil
}

func (av *Scraper) GetAcestreamId(channelId string) (string, error) {
	ch := av.Channels[channelId]
	if ch == nil {
		err := fmt.Errorf("Channel not found (id=%q)", channelId)
		return "", err
	}
	return ch.GetLink(av.client())
}

// EventByChannelAndTime returns
func (av *Scraper) EventByChannelAndTime(channel *Channel, currTime time.Time) (event *Event, live *Live) {

	for _, e := range av.Events {
		l := e.liveByChannel(channel)
		if l != nil {
			// the channel plays the event
			if event == nil {
				// it's the first event played on the channel
				event = e
				live = l
			} else {
				if currTime.Before(e.StartTime) {
					// the lives are ordered...
					return
				}
				event = e
				live = l
			}
		}
	}
	return
}

func (sch *Scraper) EventByChannel(channel *Channel) (event *Event, live *Live) {
	return sch.EventByChannelAndTime(channel, time.Now())
}

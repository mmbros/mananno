package arenavision

//package main

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/mananno/scraper"
)

type Schedule struct {
	mx         sync.RWMutex
	lastUpdate time.Time
	events     Events
}

func parseStartTime(theDay, theTime string) (time.Time, error) {
	format := "02/01/2006 15:04 MST"
	return time.Parse(format, theDay+" "+theTime)
}

func (sch *Schedule) SourceURL() string {
	return "http://arenavision.in/schedule"
}

// Get creates or updates the scheduled events from the
// "arenavision.in/schedule" page.
func (sch *Schedule) Get(client scraper.URLGetter) error {
	// get the schedule page
	resp, err := getURL(client, sch.SourceURL())
	if err != nil {
		return err
	}

	// create a goquery document from the response
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	// get the events
	events := sch.getEvents(doc)

	// updates schedule object
	sch.mx.Lock()
	sch.events = events
	// If from cache, don't update the lastUpdate timestamp
	// It uses hard coded strings not to import httpcache package.
	// XXX: Can't find a better solution.
	if sch.lastUpdate.IsZero() || resp.Header.Get("X-MMbros-Cache") != "HIT" {
		sch.lastUpdate = time.Now()
	}
	sch.mx.Unlock()

	return nil
}

func (sch *Schedule) Events() (Events, time.Time) {
	sch.mx.RLock()
	lastUpdate := sch.lastUpdate
	events := sch.events
	sch.mx.RUnlock()

	return events, lastUpdate
}

func (sch *Schedule) getEvents(doc *goquery.Document) []*Event {
	// regular expression used by getText function
	re := regexp.MustCompile(`(  +)`)

	// getText is an auxialiry function to get and clean the text
	getText := func(i int, sel *goquery.Selection) string {
		s := sel.Eq(i).Text()
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.Replace(s, "\t", " ", -1)
		s = re.ReplaceAllString(s, " ")
		s = strings.TrimSpace(s)
		return s
	}

	// find the Events
	events := []*Event{}
	doc.Find("table.auto-style1 tr").Each(func(i int, s *goquery.Selection) {
		// skip first row
		if i > 0 {
			td := s.Find("td")
			if td.Length() == 6 {
				start, err := parseStartTime(getText(0, td), getText(1, td))
				// skip last 2 rows without events
				if err == nil {

					event := &Event{StartTime: start}
					event.Sport = getText(2, td)
					event.Competition = getText(3, td)
					event.Event = getText(4, td)
					event.live = getText(5, td)
					event.Lives = stringToLives(event.live)
					// append the event
					events = append(events, event)
				}
			}
		}
	})
	return events
}

// EventByChannelAndTime returns
func (sch *Schedule) EventByChannelAndTime(channel Channel, currTime time.Time) (event *Event, live *Live) {
	events, _ := sch.Events()

	for _, e := range events {
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

func (sch *Schedule) EventByChannel(channel Channel) (event *Event, live *Live) {
	return sch.EventByChannelAndTime(channel, time.Now())
}

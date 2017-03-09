package arenavision

//package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/mananno/scraper"
)

// Channel represents an Arenavision channel.
// The values are from "1" to "30" for acestream channels.
type Channel string

// URL is the local name of the Arenavision channel.
// Examples: "av1" .. "av30"
func (ch Channel) URL() string {
	return "av" + strings.ToLower(string(ch))
}

// SourceURL is the source url of the Arenavision channel.
func (ch Channel) SourceURL() string {
	return "http://arenavision.in/av" + strings.ToLower(string(ch))
}

// FullName is the complete name of the Arenavision channel.
func (ch Channel) FullName() string {
	return "ArenaVision " + string(ch)
}

// Live represents an Arenavision live.
type Live struct {
	Channel Channel
	Lang    string
}

func (cl *Live) String() string {
	return fmt.Sprintf("{ch:%s, lang:%s}", cl.Channel, cl.Lang)
}

// Event struct
type Event struct {
	StartTime   time.Time
	Sport       string
	Competition string
	Event       string
	live        string
	Lives       []*Live
}

// Events is a slice of Event structs.
type Events []*Event

type Schedule struct {
	mx         sync.RWMutex
	lastUpdate time.Time
	events     Events
}

func (e *Event) String() string {
	return fmt.Sprintf("Event{%s, %s, %s, %s, %s}",
		e.StartTime, e.Sport, e.Competition, e.Event, e.live)
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false

}

func (e *Event) Match(filter url.Values) bool {
	sports := filter["sport"]
	if sports != nil {
		if !isValueInList(e.Sport, sports) {
			return false
		}
	}
	competitions := filter["competition"]
	if competitions != nil {
		if !isValueInList(e.Competition, competitions) {
			return false
		}
	}

	return true
}
func (events Events) FilteredBy(filter url.Values) Events {
	res := Events{}
	for _, e := range events {
		if e.Match(filter) {
			res = append(res, e)
		}
	}
	return res
}

func parseStartTime(theDay, theTime string) (time.Time, error) {
	format := "02/01/2006 15:04 MST"
	return time.Parse(format, theDay+" "+theTime)
}

func stringToLives(s string) []*Live {

	re := regexp.MustCompile(`([^\]>]+)\s+\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil
	}
	res := make([]*Live, 0, len(matches))

	for _, match := range matches {
		lang := match[2]
		channels := strings.Split(match[1], "-")
		for _, ch := range channels {
			res = append(res, &Live{
				Channel: Channel(strings.TrimSpace(ch)),
				Lang:    lang,
			})
		}
	}

	return res
}

func (sch *Schedule) SourceURL() string {
	return "http://arenavision.in/schedule"
}

func (sch *Schedule) Refresh() {

	resp, err := scraper.Get(sch.SourceURL())
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}
	events := sch.getEvents(doc)

	sch.mx.Lock()
	sch.events = events
	sch.lastUpdate = time.Now()
	sch.mx.Unlock()
}

func (sch *Schedule) Events() (Events, time.Time) {
	sch.mx.RLock()
	lastUpdate := sch.lastUpdate
	events := sch.events
	sch.mx.RUnlock()

	return events, lastUpdate
}

func (sch *Schedule) getEvents(doc *goquery.Document) []*Event {
	// getText is an auxialiry function to get and clean the text
	getText := func(i int, sel *goquery.Selection) string {
		re := regexp.MustCompile(`(  +)`)

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

/*
func main() {
	sch := new(Schedule)
	sch.Refresh()

	for j, e := range sch.Events {
		fmt.Printf("%03d) %s\n", j, e)
	}
}
*/

// returns the
func (e *Event) liveByChannel(channel Channel) *Live {
	for _, live := range e.Lives {
		if live.Channel == channel {
			return live
		}
	}
	return nil
}

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
					return
				} else {
					event = e
					live = l
				}
			}
		}
	}
	return
}

func (sch *Schedule) EventByChannel(channel Channel) (event *Event, live *Live) {
	return sch.EventByChannelAndTime(channel, time.Now())
}

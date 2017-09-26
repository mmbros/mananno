package arenavision

//package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Event struct
type Event struct {
	StartTime   time.Time
	Sport       string
	Competition string
	Event       string
	Lives       []*Live
	live        string
}

// Events is a slice of Event structs.
type Events []*Event

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

// Match check if the event matches the filter.
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

// FilteredBy returns only the events matching the filters criteria.
func (events Events) FilteredBy(filter url.Values) Events {
	res := Events{}
	for _, e := range events {
		if e.Match(filter) {
			res = append(res, e)
		}
	}
	return res
}

func (e *Event) liveByChannel(channel *Channel) *Live {
	for _, live := range e.Lives {
		if live.Channel == channel {
			return live
		}
	}
	return nil
}

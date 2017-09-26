package arenavision

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmbros/www/scraper"
)

func absURL(rel string) string {
	u, err := scraper.AbsoluteURL(baseURL, rel)
	if err != nil {
		log.Println(err)
		return ""
	}
	return u
}

func getURL(client scraper.URLGetter, URL string) (*http.Response, error) {

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	// set beget cookie
	cookie := http.Cookie{
		Name:    "beget",
		Value:   "begetok",
		Path:    "/",
		Expires: time.Now().Add(19360000000),
	}
	req.AddCookie(&cookie)

	// set request headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:52.0) Gecko/20100101 Firefox/52.0")

	return client.Do(req)
}

func parseChannels(doc *goquery.Document) Channels {
	channels := Channels{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "ArenaVision") {
			if url, _ := s.Attr("href"); len(url) > 0 {

				ch := &Channel{
					URL:  url,
					Name: text,
				}
				channels[ch.ID()] = ch
			}
		}
	})
	return channels
}

// returns "" if schedule url is not found
func parseGuideURL(doc *goquery.Document) string {
	// <a href="/iguide">EVENTS GUIDE</a>
	var url string
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Text() == "EVENTS GUIDE" {
			url, _ = s.Attr("href")
			return false // exit
		}
		return true // continue
	})
	return url
}

func parseStartTime(theDay, theTime string) (time.Time, error) {
	format := "02/01/2006 15:04 MST"
	return time.Parse(format, theDay+" "+theTime)
}

func parseEvents(doc *goquery.Document) Events {

	events := Events{}

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
					// .Lives value is not setted here
					//   event.Lives = stringToLives(event.live)
					// append the event
					events = append(events, event)
				}
			}
		}

	})
	return events
}

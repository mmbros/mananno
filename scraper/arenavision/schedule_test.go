package arenavision

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const scheduleHTML = `
</script><script src="//js.srcsmrtgs.com/js/ad.js"></script><table align="center" cellspacing="1" class="auto-style1" style="width: 100%; float: left"><tr><th class="auto-style4" style="width: 190px; height: 39px"><strong>DAY</strong></th>
<th class="auto-style4" style="width: 182px; height: 39px"><strong>TIME</strong></th>
<th class="auto-style4" style="height: 39px; width: 188px">SPORT</th>
<th class="auto-style4" style="width: 283px; height: 39px">COMPETITION</th>
<th class="auto-style4" style="width: 991px; height: 39px"><strong>EVENT</strong></th>
<th class="auto-style4" style="width: 189px; height: 39px"><strong></strong></th>
</tr><tr><td class="auto-style3" style="width: 74px">10/03/2017</td>
<td class="auto-style3" style="width: 174px">15:30 CET</td>
<td class="auto-style3" style="width: 196px">CYCLING</td>
<td class="auto-style3" style="width: 283px">UCI TOUR</td>
<td class="auto-style3" style="width: 837px">PARIS-NICE 2017. STAGE 6</td>
<td class="auto-style3" style="width: 209px">8-9 [SPA]</td>
</tr><tr><td class="auto-style3" style="width: 74px">10/03/2017</td>
<td class="auto-style3" style="width: 174px">19:00 CET</td>
<td class="auto-style3" style="width: 196px">SOCCER</td>
<td class="auto-style3" style="width: 283px">FRENCH LIGUE 1</td>
<td class="auto-style3" style="width: 837px">NICE-CAEN</td>
<td class="auto-style3" style="width: 209px">5-6 [SPA]</td>
</tr><tr><td class="auto-style3" style="width: 74px">10/03/2017</td>
<td class="auto-style3" style="width: 174px">19:00 CET</td>
<td class="auto-style3" style="width: 196px">BASKETBALL</td>
<td class="auto-style3" style="width: 283px">EUROLEAGUE</td>
<td class="auto-style3" style="width: 837px">ANADOLU EFES-BROSE BASKETS</td>
<td class="auto-style3" style="width: 209px">8-9 [SPA]</td>
</tr><tr><td class="auto-style3" style="width: 74px">10/03/2017</td>
<td class="auto-style3" style="width: 174px">20:00 CET</td>
<td class="auto-style3" style="width: 196px">SOCCER</td>
<td class="auto-style3" style="width: 283px">SPANISH LA LIGA 2</td>
<td class="auto-style3" style="width: 837px">ELCHE-ALCORCON</td>
<td class="auto-style3" style="width: 209px">1-2 [SPA]</td>
</tr><tr><td class="auto-style3" style="width: 74px">10/03/2017</td>
<td class="auto-style3" style="width: 174px">20:45 CET</td>
<td class="auto-style3" style="width: 196px">SOCCER</td>
<td class="auto-style3" style="width: 283px">SPANISH LA LIGA</td>
<td class="auto-style3" style="width: 837px">ESPANYOL-LAS PALMAS</td>
<td class="auto-style3" style="width: 209px">17-18 [SPA]</td>
</tr><tr><td class="auto-style3" style="width: 74px">10/03/2017</td>
<td class="auto-style3" style="width: 174px">20:45 CET</td>
<td class="auto-style3" style="width: 196px">SOCCER</td>
<td class="auto-style3" style="width: 283px">ITALIA SERIE A</td>
<td class="auto-style3" style="width: 837px">JUVENTUS-AC MILAN</td>
<td class="auto-style3" style="width: 209px">3-4 [SPA]<br />
		26-27 [ENG]</td>
</tr><tr><td class="auto-style8" style="width: 190px">Last update:</td>
<td class="auto-style8" style="width: 182px">10/03/2017</td>
<td class="auto-style8" style="width: 188px">14:30 CET</td>
<td class="auto-style7" style="width: 283px">TIMEZONE</td>
<td class="auto-style7" style="width: 991px">*CEST TIME - <br />
		(Madrid,Paris,Brusells)</td>
<td class="auto-style7" style="width: 189px">© Arenavision</td>
</tr></table>
`

// NewDocumentFromString returns a goquery.Document from a string.
func NewDocumentFromString(html string) (*goquery.Document, error) {
	reader := strings.NewReader(html)
	return goquery.NewDocumentFromReader(reader)
}

func TestGetSchedule(t *testing.T) {
	var expected string
	doc, err := NewDocumentFromString(scheduleHTML)
	if err != nil {
		t.Fatal(err)
	}
	sch := Schedule{}
	events := sch.getEvents(doc)

	// check number of events
	iactual := len(events)
	iexpected := 6
	if iactual != iexpected {
		t.Fatalf("Number of events: expected %d, found %d", iexpected, iactual)
	}

	// check an event
	event := events[5]
	// sport
	expected = "SOCCER"
	if event.Sport != expected {
		t.Errorf("Sport: expected %q, found %q", expected, event.Sport)
	}
	// competition
	expected = "ITALIA SERIE A"
	if event.Competition != expected {
		t.Errorf("Competition: expected %q, found %q", expected, event.Competition)
	}
	// event
	expected = "JUVENTUS-AC MILAN"
	if event.Event != expected {
		t.Errorf("Event: expected %q, found %q", expected, event.Event)
	}
	// starttime
	loc, _ := time.LoadLocation("CET")
	expectedTime := time.Date(2017, time.March, 10, 20, 45, 0, 0, loc)
	if !event.StartTime.Equal(expectedTime) {
		t.Errorf("StartTime: expected %q, found %q", expectedTime, event.StartTime)
	}
}

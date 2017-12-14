package tntvillage

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Client to the corsaro-nero sites
type Client struct{}

// Category is the category
type Category int

// Available categories
const (
	CatAll              Category = 0
	CatFilmTvEProgrammi          = 1
	CatMusica                    = 2
	CatEBooks                    = 3
	CatFilm                      = 4
	CatLinux                     = 6
	CatAnime                     = 7
	CatCartoni                   = 8
	CatMacintosh                 = 9
	CatWindowsSoftware           = 10
	CatPcGame                    = 11
	CatPlaystation               = 12
	CatStudentReleases           = 13
	CatDocumentari               = 14
	CatVideoMusicali             = 21
	CatSport                     = 22
	CatTeatro                    = 23
	CatWrestling                 = 24
	CatVarie                     = 25
	CatXBox                      = 26
	CatImmaginiESfondi           = 27
	CatAltriGiochi               = 28
	CatSerieTv                   = 29
	CatFumetteria                = 30
	CatTrash                     = 31
	CatNintendo                  = 32
	CatABook                     = 34
	CatPodcast                   = 35
	CatEdicola                   = 36
	CatMobile                    = 37
)

// SearchResults is an ordered list of search results.
type SearchResults []SearchResult

// A SearchResult contains the informations of a search result.
type SearchResult struct {
	Torrent string
	Magnet  string
	Cat     string
	Leech   string
	Seeds   string
	C       string
	Titolo  string
}

func (r *SearchResult) String() string {
	return fmt.Sprintf("%2s %3s  %3s  %s\n", r.Cat, r.Seeds, r.Leech, r.Titolo)
}

const (
	tdTorrent = iota
	tdMagnet
	tdCat
	tdLeech
	tdSeeds
	tdC
	tdTitolo
)

// doPost post a request to search the specified string and category
func (c *Client) doPost(search string, cat Category) (resp *http.Response, err error) {
	targetURL := "http://www.tntvillage.scambioetico.org/src/releaselist.php"

	form := url.Values{
		"cat":    {string(cat)},
		"page":   {"1"},
		"srcrel": {search},
	}

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "www.tntvillage.scambioetico.org")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:57.0) Gecko/20100101 Firefox/57.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "http://www.tntvillage.scambioetico.org/?releaselist")

	client := &http.Client{}

	return client.Do(req)
}

func (c *Client) doParse(doc *goquery.Document) (SearchResults, error) {
	var results SearchResults

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {

		if i > 0 { // skip first row

			var tr SearchResult

			s.Find("td").Each(func(i2 int, s2 *goquery.Selection) {
				text := strings.TrimSpace(s2.Text())
				switch i2 {
				case tdTorrent:
					tr.Torrent, _ = s2.Find("a").Attr("href")
				case tdMagnet:
					tr.Magnet, _ = s2.Find("a").Attr("href")
				case tdLeech:
					tr.Leech = text
				case tdSeeds:
					tr.Seeds = text
				case tdC:
					tr.C = text
				case tdCat:
					sCat, _ := s2.Find("a").Attr("href")
					j := strings.LastIndex(sCat, "=") + 1
					if j > 0 {
						tr.Cat = sCat[j:]
					}
				case tdTitolo:
					tr.Titolo = text
				}
			})
			results = append(results, tr)
		}
	})
	return results, nil
}

// Search the specified search string anf category
func (c *Client) Search(search string, cat Category) (SearchResults, error) {
	response, err := c.doPost(search, cat)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// check http status code
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf(response.Status)
	}

	// create a goquery document from response
	doc, err := goquery.NewDocumentFromReader(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}
	return c.doParse(doc)
}

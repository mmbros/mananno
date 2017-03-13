package ilcorsaronero

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Client to the corsaro-nero sites
type Client struct{}

// Category is the category
type Category int

// Available categories
const (
	CatAll     Category = 0
	CatFilm             = 1
	CatEBooks           = 6
	CatSerieTv          = 15
)

// SearchResults is an ordered list of search results.
type SearchResults []SearchResult

// A SearchResult contains the informations of a search result.
type SearchResult struct {
	Cat      string
	Name     string
	HREF     string
	Size     string
	Date     string
	Seeds    string
	Leech    string
	Download string
}

func (r *SearchResult) String() string {
	return fmt.Sprintf("%10s  %9s  %s  %3s  %3s  %s  %s\n", r.Cat, r.Size, r.Date, r.Seeds, r.Leech, r.Download, r.HREF)
}

const (
	tdCat = iota
	tdName
	tdSize
	tdDownload
	tdDate
	tdSeeds
	tdLeech
)

// Build a search URL with the specified search string and category
func (c *Client) buildSearchURL(search string, cat Category) string {
	// replace all `spaces` with `plus`
	query := strings.Replace(search, " ", "+", -1)
	if cat == CatAll {
		return fmt.Sprintf("http://ilcorsaronero.info/adv/%v.html", query)
	}
	return fmt.Sprintf("http://ilcorsaronero.info/adv/%v/%v.html", cat, query)
}

// Search the specified search string anf category
func (c *Client) Search(search string, cat Category) (SearchResults, error) {
	url := c.buildSearchURL(search, cat)
	log.Printf("GET %s", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results SearchResults

	doc.Find("tr[class^='odd']").Each(func(i int, s *goquery.Selection) {
		//        fmt.Printf("TD %v: %s\n", i, s.Text())
		var tr SearchResult

		s.Find("td").Each(func(i2 int, s2 *goquery.Selection) {
			text := strings.TrimSpace(s2.Text())
			switch i2 {
			case tdCat:
				tr.Cat = text
			case tdName:
				tr.HREF, _ = s2.Find("a").Attr("href")
				tr.Name = text
			case tdSize:
				tr.Size = text
			case tdDownload:
				tr.Download, _ = s2.Find("input.downarrow").Attr("value")
			case tdDate:
				tr.Date = text
			case tdSeeds:
				tr.Seeds = text
			case tdLeech:
				tr.Leech = text
			}
		})
		results = append(results, tr)
	})
	return results, nil
}

// TorrentInfo is
type TorrentInfo struct {
	HRef        string
	Magnet      string
	Cat         string
	Size        string
	Hash        string
	AnnounceURL string
	Completato  string
	Aggiunto    string
	Seeds       string
	Leech       string
}

// PrettyFormat is ...
func (ti *TorrentInfo) PrettyFormat() string {
	return fmt.Sprintf("TorrentInfo(\n"+
		"  HRef: %s\n"+
		"  Magnet: %s\n"+
		"  Cat: %s\n"+
		"  Size: %s\n"+
		"  Hash: %s\n"+
		"  AnnounceURL: %s\n"+
		"  Completato: %s\n"+
		"  Aggiunto: %s\n"+
		"  Seeds: %s\n"+
		"  Leech: %s\n"+
		")\n", ti.HRef, ti.Magnet, ti.Cat, ti.Size, ti.Hash, ti.AnnounceURL, ti.Completato, ti.Aggiunto, ti.Seeds, ti.Leech)
}

// GetTorrentInfo is ...
func (c *Client) GetTorrentInfo(href string) (*TorrentInfo, error) {
	log.Printf("corsaro.GetTorrentInfo(href=%s)", href)

	response, err := http.Get(href)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	info := TorrentInfo{}

	info.HRef = href

	info.Magnet, _ = doc.Find("a[class^='forbtn magnet']").First().Attr("href")
	//	info.cat, _ = doc.Find("td[class^='forbtn magnet']").First().Attr("href")
	/*
		<tr class="odd"><td>AnnounceURL</td>     <td><div style="width:390px;  overflow:hidden">http://tracker.tntvillage.scambioetico.org:2710/announce</div></td>  </tr>
	*/

	trimtext := func(s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	}

	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		label := tr.Find("td").First()

		switch label.Text() {
		case "Categoria bittorrent":
			/*
				<tr class="odd2">
				<td>Categoria bittorrent</td>  <td>SerieTv</td>
				</tr>
			*/
			info.Cat = trimtext(label.Next())

		case "Size":
			/*
				<tr><td>Size</td> <td>3.79 GB </td> </tr>
			*/
			info.Size = trimtext(label.Next())

		case "Hash":
			/*
				<tr class="odd2"><td>Hash</td> <td>ba9e89ac3869b03fb13100d6a0a227d62905badd</td>  </tr>
			*/
			info.Hash = trimtext(label.Next())

		case "Completato":
			/*
				<tr class="odd2"><td>Completato</td>    <td>260x</td>  </tr>
			*/
			info.Completato = trimtext(label.Next())

		case "Aggiunto":
			/*
				<tr class="odd"><td>Aggiunto</td>    <td>19.09.14  - 13:09:40</td> </tr>
			*/
			info.Aggiunto = trimtext(label.Next())

		case "AnnounceURL":
			/*
				<tr class="odd"><td>AnnounceURL</td>     <td><div style="width:390px;  overflow:hidden">http://tracker.tntvillage.scambioetico.org:2710/announce</div></td>  </tr>
			*/
			info.AnnounceURL = trimtext(tr.Find("div").First())

		case "Peers":
			/*
				<tr class="odd2">
				<td>Peers</td>
				 <td>seeds: <font color="#0066FF"> 45  </font>, leech: <font color="#006633"> 166 </font></td>
				</tr>
			*/
			font := label.Next().Find("font").First()
			info.Seeds = trimtext(font)
			info.Leech = trimtext(font.Next())

			/*
				XXX: Uploader
				XXX: Votazione
			*/

			//		default:
			//			fmt.Printf("XXX: %s\n", label.Text())
		}

	})
	fmt.Print(info.PrettyFormat())

	return &info, nil
}

// GetMagnet retrieves the magnet link from the CorsaroNero href page.
// In case of error returns the magnet link build with hash.
func (c *Client) GetMagnet(href, hash string) string {
	if ti, err := c.GetTorrentInfo(href); err == nil {
		return ti.Magnet
	}
	return "magnet:?xt=urn:btih:" + hash
}

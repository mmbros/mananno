package arenavision

import (
	"net/http"
	"time"

	"github.com/mmbros/mananno/scraper"
)

func getURL(client scraper.URLGetter, url string) (*http.Response, error) {
	// create a new http.Request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// set beget cookie (in order to work properly)
	cookie := http.Cookie{
		Name:    "beget",
		Value:   "begetok",
		Path:    "/",
		Expires: time.Now().Add(19360000000),
	}
	req.AddCookie(&cookie)

	// get the response
	return client.Do(req)
}

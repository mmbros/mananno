package scraper

import (
	"net/http"
	"net/url"
	"time"
)

// URLGetter interface defines the methods needed to get an URL.
// A standard http.Client satisfy the interface.
type URLGetter interface {
	Get(url string) (*http.Response, error)
	Do(*http.Request) (*http.Response, error)
}

// DefaultURLGetter returns an http.Client as the default URL getter.
func DefaultURLGetter() URLGetter {
	client := &http.Client{
		// Don’t use Go’s default HTTP client (in production)
		// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.q5iexu8v7
		Timeout: 90 * time.Second,
	}
	return client
}

func AbsoluteURL(base, ref string) (string, error) {
	ubase, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	uref, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	return ubase.ResolveReference(uref).String(), nil
}

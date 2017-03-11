package arenavision

import (
	"strings"

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

// GetLink get the channel source page and returns the first stream link found.
// Limits: currently it handles only acestream links.
func (ch Channel) GetLink(client scraper.URLGetter) (string, error) {
	resp, err := getURL(client, ch.SourceURL())
	if err != nil {
		return "", err
	}
	return scraper.GetFirstAcestreamLink(resp)
}

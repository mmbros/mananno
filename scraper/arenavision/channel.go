package arenavision

import "strings"

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

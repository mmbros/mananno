package arenavision

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Live represents an Arenavision live, i.e. a channel and a language.
type Live struct {
	Channel *Channel
	Lang    string
}

func (cl *Live) String() string {
	return fmt.Sprintf("{ch:%s, lang:%s}", cl.Channel, cl.Lang)
}

// stringToLives transforms a string into a slice of Live.
// Example of input string:
//    "11-12 [SPA] <br />       13-14-S3 [SRB] 14 [ITA]"
func stringToLives(s string, channels Channels) []*Live {

	re := regexp.MustCompile(`([^\]>]+)\s+\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil
	}
	res := make([]*Live, 0, len(matches))

	for _, match := range matches {
		lang := match[2]
		chids := strings.Split(match[1], "-")
		for _, chid := range chids {
			if ch := channels.Get(chid); ch == nil {
				log.Printf("Channel not found: channel id = %q", chid)
			} else {
				res = append(res, &Live{
					Channel: ch,
					Lang:    lang,
				})

			}
		}
	}

	return res
}

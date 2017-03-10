package arenavision

import (
	"fmt"
	"regexp"
	"strings"
)

// Live represents an Arenavision live, i.e. a channel and a language.
type Live struct {
	Channel Channel
	Lang    string
}

func (cl *Live) String() string {
	return fmt.Sprintf("{ch:%s, lang:%s}", cl.Channel, cl.Lang)
}

// stringToLives transforms a string into a slice of Live.
// Example of input string:
//    "11-12 [SPA] <br />       13-14-S3 [SRB] 14 [ITA]"
func stringToLives(s string) []*Live {

	re := regexp.MustCompile(`([^\]>]+)\s+\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil
	}
	res := make([]*Live, 0, len(matches))

	for _, match := range matches {
		lang := match[2]
		channels := strings.Split(match[1], "-")
		for _, ch := range channels {
			res = append(res, &Live{
				Channel: Channel(strings.TrimSpace(ch)),
				Lang:    lang,
			})
		}
	}

	return res
}

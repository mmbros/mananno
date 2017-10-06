package templates

import (
	"html/template"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"weekday": formatItWeekday,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"sum":     sum,
	"ToURL":   ToURL,
}

func ToURL(s string) template.URL {
	return template.URL(s)
}

func formatItWeekday(t time.Time) string {
	s := [...]string{"dom", "lun", "mar", "mer", "gio", "ven", "sab"}
	return s[t.Weekday()]
}

func sum(a ...int) int {
	s := 0
	for _, x := range a {
		s += x
	}
	return s
}

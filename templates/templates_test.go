package templates

import (
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	var page = PageArenavisionChannel
	wr := os.Stdout

	if err := page.Execute(wr, nil); err != nil {
		t.Error(err)
	}
}

package tntvillage

import (
	"testing"

	"github.com/mmbros/mananno/scraper"
)

func TestParse(t *testing.T) {

	doc, err := scraper.NewDocumentFromFile("the100.html")
	if err != nil {
		t.Fatal(err)
	}

	client := &Client{}

	res, err := client.doParse(doc)
	if err != nil {
		t.Fatal(err)
	}
	//	t.Log(res)
	for j, tr := range res {
		t.Log("ITEM", j, ")  ", tr)
	}

}

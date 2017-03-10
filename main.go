package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mmbros/mananno/httpcache"
	"github.com/mmbros/mananno/scraper/arenavision"
)

func main() {
	client := httpcache.NewTTL("/tmp/mananno", 5*time.Minute)

	sch := new(arenavision.Schedule)
	err := sch.Get(client)
	if err != nil {
		log.Fatal(err)
	}

	events, lastUpdate := sch.Events()
	fmt.Printf("Events %d - LastUpdate %s\n", len(events), lastUpdate)
	for i, event := range events {
		fmt.Printf("%2d) %v\n", i, event)
	}

}

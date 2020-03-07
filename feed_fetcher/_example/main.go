package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/feed_fetcher"
)

func main() {
	url := os.Args[1]
	fp := gofeed.NewParser()
	fetcher := feed_fetcher.NewFeedParser(fp, &http.Client{})
	start := time.Now()
	feed, err := fetcher.Fetch(url)
	log.Println("time: ", time.Now().Sub(start))
	if err != nil {
		log.Fatalln(err)
	}
	x, _ := json.MarshalIndent(feed.Feed(), "", "\t")
	log.Println("==== feed =====")
	log.Println(string(x))
	x, _ = json.MarshalIndent(feed.Entries(), "", "\t")
	log.Println("=== entries ===")
	log.Println(string(x))
}

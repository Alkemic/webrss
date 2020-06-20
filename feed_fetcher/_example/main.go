package main

import (
	"context"
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
	ctx := context.Background()
	feed, err := fetcher.Fetch(ctx, url)
	log.Println("time: ", time.Now().Sub(start))
	if err != nil {
		log.Fatalln(err)
	}
	x, _ := json.MarshalIndent(feed.Feed(ctx), "", "\t")
	log.Println("==== feed =====")
	log.Println(string(x))
	x, _ = json.MarshalIndent(feed.Entries(ctx), "", "\t")
	log.Println("=== entries ===")
	log.Println(string(x))
}

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Alkemic/webrss/updater"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/config"
	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

var (
	configFile = flag.String("config", "config.yml", "")
	logger     = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Ldate)
)

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		logger.Fatalln("cannot load config file: ", err)
	}

	db, err, closeFn := initDB(cfg)
	if err != nil {
		logger.Fatalf("cannot instantiate db: %s", err)
	}
	defer closeFn()
	fp := gofeed.NewParser()
	feedFetcher := feed_fetcher.NewFeedParser(fp, &http.Client{})

	feedRepository := repository.NewFeedRepository(db)
	entryRepository := repository.NewEntryRepository(db, cfg.PerPage)
	feedService := webrss.NewFeedService(feedRepository, entryRepository, feedFetcher)
	updateService := updater.New(feedRepository, feedService, feedFetcher, logger)
	if err := updateService.Run(context.Background()); err != nil {
		logger.Println("got error updating feeds:", err)
		return
	}
	logger.Println("done.")
}

func initDB(cfg *config.Config) (*sqlx.DB, error, func()) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open connection to db: %w", err), nil
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err), nil
	}
	return sqlx.NewDb(db, "mysql"), nil, func() {
		db.Close()
	}
}

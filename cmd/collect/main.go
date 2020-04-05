package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/config"
	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/updater"
	"github.com/Alkemic/webrss/webrss"
)

func main() {
	flag.Parse()
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Ldate)
	cfg := config.LoadConfig()
	db, err, closeFn := initDB(cfg)
	if err != nil {
		logger.Fatalf("cannot instantiate db: %s", err)
	}
	defer closeFn()
	fp := gofeed.NewParser()
	httpClient := &http.Client{}
	feedFetcher := feed_fetcher.NewFeedParser(fp, httpClient)

	feedRepository := repository.NewFeedRepository(db)
	entryRepository := repository.NewEntryRepository(db, cfg.PerPage)
	transactionRepository := repository.NewTransactionRepository(db)
	webrssService := webrss.NewService(logger, nil, feedRepository, entryRepository, transactionRepository, httpClient, feedFetcher)
	updateService := updater.New(feedRepository, webrssService, feedFetcher, logger)
	if err := updateService.Run(context.Background()); err != nil {
		logger.Println("got error updating feeds:", err)
		return
	}
	logger.Println("done.")
}

func initDB(cfg *config.Config) (*sqlx.DB, error, func()) {
	db, err := sql.Open("mysql", cfg.DBDSN)
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

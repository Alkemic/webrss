package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/Alkemic/webrss/updater"

	"github.com/Alkemic/go-route"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/config"
	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/handler"
	"github.com/Alkemic/webrss/repository"
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

	categoryRepository := repository.NewCategoryRepository(db)
	feedRepository := repository.NewFeedRepository(db)
	entryRepository := repository.NewEntryRepository(db, cfg.PerPage)
	transactionRepository := repository.NewTransactionRepository(db)
	webrssService := webrss.NewService(logger, categoryRepository, feedRepository, entryRepository, transactionRepository, httpClient, feedFetcher)
	categoryHandler := handler.NewCategory(logger, webrssService)
	entryHandler := handler.NewEntry(logger, webrssService)
	feedHandler := handler.NewFeed(logger, webrssService)
	updateService := updater.New(feedRepository, webrssService, feedFetcher, logger)
	app := webrss.New(logger, cfg, categoryHandler, feedHandler, entryHandler, updateService, time.Hour)
	app.AddOnExit(closeFn)
	go app.Updater(context.Background())
	if err := app.Run(); err != nil {
		logger.Fatalln("application exited with error: ", err)
	}
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

func panicDefer(rw http.ResponseWriter, req *http.Request, logger *log.Logger) {
	if r := recover(); r != nil {
		if logger != nil {
			logger.Printf("Panic occured:\n%s\n", r)
			logger.Println("stacktrace:\n" + string(debug.Stack()))
			logger.Println("Panic end.")
		} else {
			log.Printf("Panic occured:\n%s\n", r)
			log.Println("stacktrace:\n" + string(debug.Stack()))
			log.Println("Panic end.")
		}
		route.InternalServerError(rw, req)
	}
}

func PanicInterceptorWithLogger(logger *log.Logger) func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			defer panicDefer(rw, req, logger)

			f(rw, req)
		}
	}
}

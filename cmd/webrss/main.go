package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/account"
	"github.com/Alkemic/webrss/config"
	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/handler"
	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/updater"
	"github.com/Alkemic/webrss/webrss"
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Llongfile|log.Ldate)
	cfg := config.LoadConfig()
	db, err, closeFn := initDB(cfg)
	if err != nil {
		logger.Fatalf("cannot instantiate db: %s", err)
	}
	defer closeFn()

	fp := gofeed.NewParser()
	httpClient := &http.Client{}
	feedFetcher := feed_fetcher.NewFeedParser(fp, httpClient)

	//userRepository := repository.NewUserRepository(db)
	settingsRepository := repository.NewSettingsRepository(db)
	sessionRepository := repository.NewSessionRepository(28 * 24 * time.Hour)
	authenticateHandler := account.NewAuthenticateHandler(logger, settingsRepository, sessionRepository)
	authenticateMiddleware := account.NewAuthenticateMiddleware(logger, settingsRepository, sessionRepository)

	categoryRepository := repository.NewCategoryRepository(db)
	feedRepository := repository.NewFeedRepository(db)
	entryRepository := repository.NewEntryRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	webrssService := webrss.NewService(logger, categoryRepository, feedRepository, entryRepository, transactionRepository, httpClient, feedFetcher)
	categoryHandler := handler.NewCategory(logger, webrssService)
	entryHandler := handler.NewEntry(logger, webrssService, cfg.PerPage)
	feedHandler := handler.NewFeed(logger, webrssService)
	updateService := updater.New(feedRepository, webrssService, feedFetcher, logger)
	app := webrss.New(logger, cfg, categoryHandler, feedHandler, entryHandler, authenticateHandler, authenticateMiddleware, updateService, time.Hour)
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

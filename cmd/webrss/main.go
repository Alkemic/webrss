package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/Alkemic/webrss/config"
	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
	"github.com/Alkemic/webrss/webrss/category"
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

	categoryRepository := repository.NewCategoryRepository(db)
	feedRepository := repository.NewFeedRepository(db)
	columnService := webrss.NewCategoryService(categoryRepository, feedRepository)
	categoryHandler := category.NewHandler(columnService, logger)

	routes := route.RegexpRouter{}
	routes.Add("^/api/category", categoryHandler.GetRoutes())
	routes.Add("^/favicon.ico$", favicon)
	routes.Add("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("webrss/static"))))
	routes.Add("/", index)

	//handler := middleware.TimeTrack(logger)(middleware.PanicInterceptorWithLogger(logger)(routes.ServeHTTP))
	handler := middleware.TimeTrack(logger)(PanicInterceptorWithLogger(logger)(routes.ServeHTTP))
	bindAddr := fmt.Sprintf("%s:%d", cfg.Run.Host, cfg.Run.Port)
	if err := http.ListenAndServe(bindAddr, handler); err != nil {
		log.Fatalln("exited with error: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "webrss/templates/index.html")
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "webrss/static/images/favicon.ico")
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

package webrss

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"

	"github.com/Alkemic/webrss/config"
)

type handler interface {
	GetRoutes() route.RegexpRouter
}

type feedsUpdater interface {
	Run(ctx context.Context) error
}

type App struct {
	logger          *log.Logger
	cfg             *config.Config
	routes          route.RegexpRouter
	categoryHandler handler
	feedHandler     handler
	entryHandler    handler
	feedsUpdater    feedsUpdater
	updaterInterval time.Duration

	onExit []func()
}

func New(logger *log.Logger, cfg *config.Config, categoryHandler handler, feedHandler handler, entryHandler handler, feedsUpdater feedsUpdater, updaterInterval time.Duration) App {
	app := App{
		logger:          logger,
		cfg:             cfg,
		categoryHandler: categoryHandler,
		feedHandler:     feedHandler,
		entryHandler:    entryHandler,
		feedsUpdater:    feedsUpdater,
		updaterInterval: updaterInterval,
	}
	app.routes.Add("^/api/category", categoryHandler.GetRoutes())
	app.routes.Add("^/api/entry", entryHandler.GetRoutes())
	app.routes.Add("^/api/feed", feedHandler.GetRoutes())
	app.routes.Add("^/favicon.ico$", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})
	app.routes.Add("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	app.routes.Add("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	return app
}

func (a App) Run() error {
	defer a.execOnExit()
	handler := middleware.TimeTrack(a.logger)(middleware.PanicInterceptorWithLogger(a.logger)(a.routes.ServeHTTP))
	if err := http.ListenAndServe(a.cfg.BindAdr, handler); err != http.ErrServerClosed {
		return fmt.Errorf("exited with error: %w", err)
	}
	return nil
}

func (a *App) AddOnExit(fn func()) {
	a.onExit = append(a.onExit, fn)
}

func (a *App) execOnExit() {
	for _, fn := range a.onExit {
		fn()
	}
}

func (a App) Updater(ctx context.Context) {
	if !a.cfg.RunUpdater {
		a.logger.Println("updater won't be running")
		return
	}
	ticker := time.NewTicker(a.updaterInterval)
	for {
		if err := a.feedsUpdater.Run(ctx); err != nil {
			a.logger.Println("task returned an error: ", err)
		}
		<-ticker.C
	}
}

package webrss

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alkemic/go-route/middleware"

	"github.com/Alkemic/webrss/config"

	"github.com/Alkemic/go-route"
)

type handler interface {
	GetRoutes() route.RegexpRouter
}

type App struct {
	logger          *log.Logger
	cfg             *config.Config
	routes          route.RegexpRouter
	categoryHandler handler
	feedHandler     handler
	entryHandler    handler

	onExit []func()
}

func New(logger *log.Logger, cfg *config.Config, categoryHandler handler, feedHandler handler, entryHandler handler) App {
	app := App{
		logger:          logger,
		cfg:             cfg,
		categoryHandler: categoryHandler,
		feedHandler:     feedHandler,
		entryHandler:    entryHandler,
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
	bindAddr := fmt.Sprintf("%s:%d", a.cfg.Run.Host, a.cfg.Run.Port)

	if err := http.ListenAndServe(bindAddr, handler); err != http.ErrServerClosed {
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

//func index(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "templates/index.html")
//}

//func favicon(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "static/images/favicon.ico")
//}

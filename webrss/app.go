package webrss

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"

	"github.com/Alkemic/webrss/account"
	"github.com/Alkemic/webrss/config"
	"github.com/Alkemic/webrss/repository"
)

type handler interface {
	GetRoutes() *route.RegexpRouter
}

type feedsUpdater interface {
	Run(ctx context.Context) error
}

type App struct {
	logger          *log.Logger
	cfg             *config.Config
	routes          *route.RegexpRouter
	categoryHandler handler
	feedHandler     handler
	entryHandler    handler
	feedsUpdater    feedsUpdater
	updaterInterval time.Duration

	onExit []func()
}

func New(logger *log.Logger, cfg *config.Config, categoryHandler handler, feedHandler handler, entryHandler handler,
	authenticateHandler *account.AuthenticateHandler, authenticateMiddleware *account.Middleware, feedsUpdater feedsUpdater, updaterInterval time.Duration) App {
	app := App{
		logger:          logger,
		cfg:             cfg,
		routes:          route.New(),
		categoryHandler: categoryHandler,
		feedHandler:     feedHandler,
		entryHandler:    entryHandler,
		feedsUpdater:    feedsUpdater,
		updaterInterval: updaterInterval,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	app.routes.Add(account.LoginPageURL, authenticateHandler.Login)
	app.routes.Add(account.LogoutPageURL, authenticateHandler.Logout)

	app.routes.Add("^/api/category", categoryHandler.GetRoutes().AddMiddleware(authenticateMiddleware.LoginRequiredMiddleware))
	app.routes.Add("^/api/entry", entryHandler.GetRoutes().AddMiddleware(authenticateMiddleware.LoginRequiredMiddleware))
	app.routes.Add("^/api/feed", feedHandler.GetRoutes().AddMiddleware(authenticateMiddleware.LoginRequiredMiddleware))
	app.routes.Add("^/favicon.ico$", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})
	app.routes.Add("^/api/user/$", setHeaders(authenticateMiddleware.LoginRequiredMiddleware(authenticateHandler.Edit)))
	app.routes.Add("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	app.routes.Add("/", authenticateMiddleware.LoginRequiredMiddleware(func(rw http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.New("index.html").Funcs(map[string]interface{}{
			"marshal": func(v interface{}) template.JS {
				a, _ := json.Marshal(v)
				return template.JS(a)
			},
		}).Delims("[[", "]]").ParseFiles("templates/index.html"))
		tmplData := struct {
			User      repository.User
			LogoutURL string
		}{
			User:      account.GetUser(req),
			LogoutURL: account.LogoutPageURL,
		}
		if err := tmpl.ExecuteTemplate(rw, "index.html", tmplData); err != nil {
			logger.Println("cannot execute template:", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}))

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

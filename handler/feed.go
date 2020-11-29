package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"
	"gopkg.in/go-playground/validator.v9"

	httphelper "github.com/Alkemic/webrss/http"
	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type FeedValid struct {
	FeedURL        string `validate:"required,min=3,max=255,url" json:"feed_url"`
	FeedFaviconURL string `validate:"max=255" json:"site_favicon_url"`
	FeedTitle      string `validate:"max=255" json:"feed_title"`
	Category       int64  `validate:"required"`
}

type feedHandler struct {
	logger        *log.Logger
	webrssService webrssService
}

func NewFeed(logger *log.Logger, service webrssService) *feedHandler {
	return &feedHandler{
		webrssService: service,
		logger:        logger,
	}
}

func (h *feedHandler) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Println("error reading body:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	feedData := FeedValid{}
	if err := json.Unmarshal(body, &feedData); err != nil {
		h.logger.Println("can't unmarshal body:", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = validator.New().Struct(feedData); err != nil {
		h.logger.Println("validation error:", err)
		http.Error(rw, "validation error", http.StatusBadRequest)
		return
	}

	if err := h.webrssService.CreateFeed(req.Context(), feedData.FeedURL, feedData.Category); err != nil {
		h.logger.Println("error creating feed:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *feedHandler) Update(rw http.ResponseWriter, req *http.Request) {
	id, err := httphelper.GetIntParam(req, "id")
	if err != nil {
		h.logger.Println("cannot get param 'id': ", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Println("error reading body: ", err)
		http.Error(rw, "can't read body", http.StatusBadRequest)
		return
	}
	feedData := FeedValid{}
	if err := json.Unmarshal(body, &feedData); err != nil {
		h.logger.Println("can't unmarshal body:", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = validator.New().Struct(feedData); err != nil {
		h.logger.Println("validation error:", err)
		http.Error(rw, "validation error", http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	feed, err := h.webrssService.GetFeed(ctx, id)
	if err != nil {
		h.logger.Println("error getting feed: ", err)
		http.Error(rw, "error getting feed", http.StatusInternalServerError)
		return
	}
	feed.FeedUrl = feedData.FeedURL
	feed.FeedTitle = feedData.FeedTitle
	feed.CategoryID = feedData.Category
	feed.SiteFaviconUrl = repository.NewNullString(feedData.FeedFaviconURL)

	if err := h.webrssService.UpdateFeed(ctx, feed); err != nil {
		h.logger.Println("error updating category: ", err)
		http.Error(rw, "error updating category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *feedHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	idRaw, ok := route.GetParam(req, "id")
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		h.logger.Println("cannot convert param 'id' to int: ", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	feed, err := h.webrssService.GetFeed(ctx, int64(id))
	if err != nil {
		h.logger.Println("cannot get feed: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := h.webrssService.DeleteFeed(ctx, feed); err != nil {
		h.logger.Println("cannot delete feed: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (r *feedHandler) GetRoutes() *route.RegexpRouter {
	resource := webrss.RESTEndPoint{
		Delete: r.Delete,
		Put:    r.Update,
	}
	collection := webrss.RESTEndPoint{
		Post: r.Create,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	routing := route.New()
	routing.Add(`^/?$`, setHeaders(collection.Dispatch))
	routing.Add(`^/(?P<id>\d+)/?$`, setHeaders(resource.Dispatch))

	return routing
}

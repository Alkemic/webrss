package entry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"

	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type entryService interface {
	Get(ctx context.Context, id int64) (repository.Entry, error)
	ListForFeed(ctx context.Context, feedID, page int64) ([]repository.Entry, error)
}

type restHandler struct {
	entryService entryService
	logger       *log.Logger
}

func NewHandler(entryService entryService, logger *log.Logger) *restHandler {
	return &restHandler{
		entryService: entryService,
		logger:       logger,
	}
}

func (h *restHandler) Get(rw http.ResponseWriter, req *http.Request) {
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

	entry, err := h.entryService.Get(req.Context(), int64(id))
	if err != nil {
		h.logger.Println("error getting entry: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(entry); err != nil {
		h.logger.Println("cannot serialize entries: ", err)
	}
}

func (h *restHandler) List(rw http.ResponseWriter, req *http.Request) {
	feedID, _, err := getIntParam("feed", req)
	if err != nil {
		h.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	page, ok, err := getIntParam("page", req)
	if err != nil && ok {
		h.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if page < 1 {
		page = 1
	}

	entries, err := h.entryService.ListForFeed(req.Context(), int64(feedID), int64(page))
	if err != nil {
		h.logger.Println("cannot fetch entries: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"objects": entries,
	}
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		h.logger.Println("cannot serialize entries: ", err)
	}
}

func getIntParam(key string, req *http.Request) (int, bool, error) {
	query := req.URL.Query()
	rawValue := query.Get(key)
	if rawValue == "" {
		return 0, false, fmt.Errorf("missing '%s' param", key)
	}
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, true, fmt.Errorf("'%s' is not int: %w", key, err)
	}
	return value, true, nil
}

func (r *restHandler) GetRoutes() route.RegexpRouter {
	resource := webrss.RESTEndPoint{
		Get: r.Get,
	}
	collection := webrss.RESTEndPoint{
		Get: r.List,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	routing := route.RegexpRouter{}
	routing.Add(`^/?$`, setHeaders(collection.Dispatch))
	routing.Add(`^/(?P<id>\d+)/?$`, setHeaders(resource.Dispatch))

	return routing
}

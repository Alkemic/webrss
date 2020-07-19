package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"

	"github.com/Alkemic/webrss/webrss"
)

type entryHandler struct {
	logger        *log.Logger
	webrssService webrssService
	perPage       int
}

func NewEntry(logger *log.Logger, webrssService webrssService, perPage int) *entryHandler {
	return &entryHandler{
		logger:        logger,
		webrssService: webrssService,
		perPage:       perPage,
	}
}

func (h *entryHandler) Get(rw http.ResponseWriter, req *http.Request) {
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

	entry, err := h.webrssService.GetEntry(req.Context(), int64(id))
	if err != nil {
		h.logger.Println("error getting entry: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(rw).Encode(entry); err != nil {
		h.logger.Println("cannot serialize entries: ", err)
	}
}

func getPage(req *http.Request) (int64, error) {
	page, ok, err := getIntParam("page", req)
	if err != nil && ok {
		return 0, err
	}
	if page < 1 {
		page = 1
	}
	return int64(page), nil
}

func (h *entryHandler) Search(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	phrase := query.Get("phrase")
	if phrase == "" {
		h.logger.Println("missing phrase in request")
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	page, err := getPage(req)
	if err != nil {
		h.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	entries, err := h.webrssService.Search(req.Context(), phrase, page, h.perPage)
	if err != nil {
		h.logger.Println("cannot fetch entries: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	nextPage := ""
	if len(entries) == h.perPage {
		nextPage = fmt.Sprintf("/api/entry/search/?phrase=%s&page=%d", phrase, page+1)
	}
	data := map[string]interface{}{
		"objects": entries,
		"meta":    map[string]string{"next": nextPage},
	}

	if err := json.NewEncoder(rw).Encode(data); err != nil {
		h.logger.Println("cannot serialize entries: ", err)
	}
}

func (h *entryHandler) List(rw http.ResponseWriter, req *http.Request) {
	feedID, _, err := getIntParam("feed", req)
	if err != nil {
		h.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	page, err := getPage(req)
	if err != nil {
		h.logger.Println(err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	entries, err := h.webrssService.ListEntriesForFeed(req.Context(), int64(feedID), page, h.perPage)
	if err != nil {
		h.logger.Println("cannot fetch entries: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	nextPage := ""
	if len(entries) == h.perPage {
		nextPage = fmt.Sprintf("/api/entry/?feed=%d&page=%d", feedID, page+1)
	}
	data := map[string]interface{}{
		"objects": entries,
		"meta":    map[string]string{"next": nextPage},
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

func (r *entryHandler) GetRoutes() route.RegexpRouter {
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
	routing.Add(`^/search/?$`, setHeaders(r.Search))
	routing.Add(`^/(?P<id>\d+)/?$`, setHeaders(resource.Dispatch))

	return routing
}

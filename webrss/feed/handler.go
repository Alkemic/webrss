package feed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"
	"gopkg.in/go-playground/validator.v9"

	"github.com/Alkemic/webrss/webrss"
)

type FeedValid struct {
	FeedURL  string `validate:"required,min=3,max=255" json:"feed_url"`
	Category int64  `validate:"required"`
}

type feedService interface {
	Create(feedURL string, categoryID int64) error
}

type restHandler struct {
	logger      *log.Logger
	feedService feedService
}

func NewHandler(logger *log.Logger, feedService feedService) *restHandler {
	return &restHandler{
		feedService: feedService,
		logger:      logger,
	}
}

func (h *restHandler) Create(rw http.ResponseWriter, req *http.Request) {
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
	log.Println(string(body))

	if err := h.feedService.Create(feedData.FeedURL, feedData.Category); err != nil {
		h.logger.Println("error creating feed:", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (r *restHandler) GetRoutes() route.RegexpRouter {
	resource := webrss.RESTEndPoint{}
	collection := webrss.RESTEndPoint{
		Post: r.Create,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	routing := route.RegexpRouter{}
	routing.Add(`^/?$`, setHeaders(collection.Dispatch))
	routing.Add(`^/(?P<id>\d+)/?$`, setHeaders(resource.Dispatch))

	return routing
}

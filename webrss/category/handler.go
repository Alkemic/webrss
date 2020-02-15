package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Alkemic/go-route/middleware"

	"github.com/Alkemic/go-route"

	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type categoryService interface {
	Get(id int64) (repository.Category, error)
	List(params ...string) ([]repository.Category, error)
	Delete(id int64) error
	Update(repository.Category) error
}

type restHandler struct {
	categoryService categoryService
	logger          *log.Logger
}

func NewHandler(categoryService categoryService, logger *log.Logger) *restHandler {
	return &restHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

func (h *restHandler) List(rw http.ResponseWriter, req *http.Request) {
	categories, err := h.categoryService.List()
	if err != nil {
		h.logger.Println("cannot fetch categories: ", err)
		return
	}
	data := map[string]interface{}{
		"objects": categories,
	}
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		h.logger.Println("cannot serialize categories: ", err)
	}
}

func (h *restHandler) Create(rw http.ResponseWriter, req *http.Request) {}

func (h *restHandler) Get(rw http.ResponseWriter, req *http.Request)    {}
func (h *restHandler) Delete(rw http.ResponseWriter, req *http.Request) {}

func (h *restHandler) GetRoutes() route.RegexpRouter {
	categoryResource := webrss.RESTEndPoint{
		Get:    h.Get,
		Delete: h.Delete,
	}
	categoryCollection := webrss.RESTEndPoint{
		Get: h.List,
	}

	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	setHeaders := middleware.SetHeaders(headers)

	categoryRouting := route.RegexpRouter{}
	categoryRouting.Add(`^/?$`, setHeaders(categoryCollection.Dispatch))
	categoryRouting.Add(`^/(?P<id>\d+)/$`, setHeaders(categoryResource.Dispatch))

	return categoryRouting
}

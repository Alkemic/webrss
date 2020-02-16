package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"

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

func (h *restHandler) GetRoutes() route.RegexpRouter {
	resource := webrss.RESTEndPoint{}
	collection := webrss.RESTEndPoint{
		Get: h.List,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	routing := route.RegexpRouter{}
	routing.Add(`^/?$`, setHeaders(collection.Dispatch))
	routing.Add(`^/(?P<id>\d+)/$`, setHeaders(resource.Dispatch))

	return routing
}

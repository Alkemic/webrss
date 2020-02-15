package category

import (
	"log"
	"net/http"

	"github.com/Alkemic/go-route"

	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type entryService interface {
	Get(id int64) (repository.Category, error)
	List(params ...string) ([]repository.Category, error)
	Delete(id int64) error
	MoveUp(id int64) error
	MoveDown(id int64) error
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

func (r *restHandler) List(rw http.ResponseWriter, req *http.Request)   {}
func (r *restHandler) Create(rw http.ResponseWriter, req *http.Request) {}

func (r *restHandler) Get(rw http.ResponseWriter, req *http.Request)    {}
func (r *restHandler) Delete(rw http.ResponseWriter, req *http.Request) {}

func (r *restHandler) GetRoutes() route.RegexpRouter {
	categoryResource := webrss.RESTEndPoint{
		Get:    r.Get,
		Delete: r.Delete,
	}
	categoryCollection := webrss.RESTEndPoint{
		Get: r.List,
	}

	categoryRouting := route.RegexpRouter{}
	categoryRouting.Add(`^/entry/?$`, categoryCollection.Dispatch)
	categoryRouting.Add(`^/entry/(?P<id>\d+)/$`, categoryResource.Dispatch)

	return categoryRouting
}

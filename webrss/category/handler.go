package category

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"
	"gopkg.in/go-playground/validator.v9"

	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type categoryService interface {
	Get(id int64) (repository.Category, error)
	List(params ...string) ([]repository.Category, error)
	Delete(id int64) error
	Update(repository.Category) error
	Create(repository.Category) error
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
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"objects": categories,
	}
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		h.logger.Println("cannot serialize categories: ", err)
	}
}

type Category struct {
	Title string `validate:"required,min=3,max=255"`
}

func (h *restHandler) Create(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Println("error reading body: ", err)
		http.Error(rw, "can't read body", http.StatusBadRequest)
		return
	}
	newCategory := Category{}
	if err := json.Unmarshal(body, &newCategory); err != nil {
		h.logger.Println("can't unmarshal body: ", err)
		http.Error(rw, "can't unmarshal body", http.StatusBadRequest)
		return
	}
	if err = validator.New().Struct(newCategory); err != nil {
		h.logger.Println("validation error: ", err)
		http.Error(rw, "validation error", http.StatusBadRequest)
		return
	}
	log.Println(string(body))

	category := repository.Category{
		Title: newCategory.Title,
	}
	if err := h.categoryService.Create(category); err != nil {
		h.logger.Println("error creating category: ", err)
		http.Error(rw, "error creating category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *restHandler) GetRoutes() route.RegexpRouter {
	resource := webrss.RESTEndPoint{}
	collection := webrss.RESTEndPoint{
		Get:  h.List,
		Post: h.Create,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	routing := route.RegexpRouter{}
	routing.Add(`^/?$`, setHeaders(collection.Dispatch))
	routing.Add(`^/(?P<id>\d+)/$`, setHeaders(resource.Dispatch))

	return routing
}

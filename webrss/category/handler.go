package category

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"
	"gopkg.in/go-playground/validator.v9"

	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type categoryService interface {
	Get(ctx context.Context, id int64) (repository.Category, error)
	List(ctx context.Context, params ...string) ([]repository.Category, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, category repository.Category) error
	Create(ctx context.Context, category repository.Category) error
	MoveUp(ctx context.Context, id int64) error
	MoveDown(ctx context.Context, id int64) error
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
	categories, err := h.categoryService.List(req.Context())
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
	if err := h.categoryService.Create(req.Context(), category); err != nil {
		h.logger.Println("error creating category: ", err)
		http.Error(rw, "error creating category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *restHandler) MoveUp(rw http.ResponseWriter, req *http.Request) {
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
	if err := h.categoryService.MoveUp(req.Context(), int64(id)); err != nil {
		h.logger.Println("error moving up: ", err)
		http.Error(rw, "error moving up", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *restHandler) MoveDown(rw http.ResponseWriter, req *http.Request) {
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
	if err := h.categoryService.MoveDown(req.Context(), int64(id)); err != nil {
		h.logger.Println("error moving up: ", err)
		http.Error(rw, "error moving up", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *restHandler) Delete(rw http.ResponseWriter, req *http.Request) {
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
	if err := h.categoryService.Delete(req.Context(), int64(id)); err != nil {
		h.logger.Println("error deleting category: ", err)
		http.Error(rw, "error deleting category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *restHandler) GetRoutes() route.RegexpRouter {
	resource := webrss.RESTEndPoint{
		Delete: h.Delete,
	}
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
	routing.Add(`^/(?P<id>\d+)/move_up$`, setHeaders(middleware.AllowedMethods([]string{http.MethodPost})(h.MoveUp)))
	routing.Add(`^/(?P<id>\d+)/move_down$`, setHeaders(middleware.AllowedMethods([]string{http.MethodPost})(h.MoveDown)))

	return routing
}

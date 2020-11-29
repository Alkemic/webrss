package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Alkemic/go-route"
	"github.com/Alkemic/go-route/middleware"
	"gopkg.in/go-playground/validator.v9"

	httphelper "github.com/Alkemic/webrss/http"
	"github.com/Alkemic/webrss/repository"
	"github.com/Alkemic/webrss/webrss"
)

type webrssService interface {
	GetCategory(ctx context.Context, id int64) (repository.Category, error)
	ListCategories(ctx context.Context, params ...string) ([]repository.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	UpdateCategory(ctx context.Context, category repository.Category) error
	CreateCategory(ctx context.Context, category repository.Category) error
	MoveCategoryUp(ctx context.Context, id int64) error
	MoveCategoryDown(ctx context.Context, id int64) error

	GetFeed(ctx context.Context, id int64) (repository.Feed, error)
	CreateFeed(ctx context.Context, feedURL string, categoryID int64) error
	DeleteFeed(ctx context.Context, feed repository.Feed) error
	UpdateFeed(ctx context.Context, feed repository.Feed) error

	SaveEntries(ctx context.Context, feedID int64, entries []repository.Entry) error

	GetEntry(ctx context.Context, id int64) (repository.Entry, error)
	Search(ctx context.Context, phrase string, page int64, perPage int) ([]repository.Entry, error)
	ListEntriesForFeed(ctx context.Context, feedID, page int64, perPage int) ([]repository.Entry, error)
}

type categoryHandler struct {
	webrssService webrssService
	logger        *log.Logger
}

func NewCategory(logger *log.Logger, categoryService webrssService) *categoryHandler {
	return &categoryHandler{
		logger:        logger,
		webrssService: categoryService,
	}
}

func (h *categoryHandler) List(rw http.ResponseWriter, req *http.Request) {
	categories, err := h.webrssService.ListCategories(req.Context())
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

func (h *categoryHandler) Create(rw http.ResponseWriter, req *http.Request) {
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
	if err := h.webrssService.CreateCategory(req.Context(), category); err != nil {
		h.logger.Println("error creating category: ", err)
		http.Error(rw, "error creating category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *categoryHandler) Update(rw http.ResponseWriter, req *http.Request) {
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
	upCategory := Category{}
	if err := json.Unmarshal(body, &upCategory); err != nil {
		h.logger.Println("can't unmarshal body: ", err)
		http.Error(rw, "can't unmarshal body", http.StatusBadRequest)
		return
	}
	if err = validator.New().Struct(upCategory); err != nil {
		h.logger.Println("validation error: ", err)
		http.Error(rw, "validation error", http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	category, err := h.webrssService.GetCategory(ctx, id)
	if err != nil {
		h.logger.Println("error getting category: ", err)
		http.Error(rw, "error getting category", http.StatusInternalServerError)
		return
	}
	category.Title = upCategory.Title
	if err := h.webrssService.UpdateCategory(ctx, category); err != nil {
		h.logger.Println("error updating category: ", err)
		http.Error(rw, "error updating category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *categoryHandler) MoveUp(rw http.ResponseWriter, req *http.Request) {
	id, err := httphelper.GetIntParam(req, "id")
	if err != nil {
		h.logger.Println("cannot get param 'id': ", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.webrssService.MoveCategoryUp(req.Context(), id); err != nil {
		h.logger.Println("error moving up: ", err)
		http.Error(rw, "error moving up", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *categoryHandler) MoveDown(rw http.ResponseWriter, req *http.Request) {
	id, err := httphelper.GetIntParam(req, "id")
	if err != nil {
		h.logger.Println("cannot get param 'id': ", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.webrssService.MoveCategoryDown(req.Context(), id); err != nil {
		h.logger.Println("error moving up: ", err)
		http.Error(rw, "error moving up", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *categoryHandler) Delete(rw http.ResponseWriter, req *http.Request) {
	id, err := httphelper.GetIntParam(req, "id")
	if err != nil {
		h.logger.Println("cannot get param 'id': ", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.webrssService.DeleteCategory(req.Context(), id); err != nil {
		h.logger.Println("error deleting category: ", err)
		http.Error(rw, "error deleting category", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, `{"status":"ok"}`)
}

func (h *categoryHandler) GetRoutes() *route.RegexpRouter {
	resource := webrss.RESTEndPoint{
		Delete: h.Delete,
		Post:   h.Update,
	}
	collection := webrss.RESTEndPoint{
		Get:  h.List,
		Post: h.Create,
	}

	setHeaders := middleware.SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})

	routing := route.New()
	routing.Add(`^/?$`, setHeaders(collection.Dispatch))
	routing.Add(`^/(?P<id>\d+)/$`, setHeaders(resource.Dispatch))
	routing.Add(`^/(?P<id>\d+)/move_up$`, setHeaders(middleware.AllowedMethods([]string{http.MethodPost})(h.MoveUp)))
	routing.Add(`^/(?P<id>\d+)/move_down$`, setHeaders(middleware.AllowedMethods([]string{http.MethodPost})(h.MoveDown)))

	return routing
}

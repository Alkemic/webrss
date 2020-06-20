package webrss

import "net/http"

type RESTEndPoint struct {
	Get     http.HandlerFunc
	Post    http.HandlerFunc
	Put     http.HandlerFunc
	Delete  http.HandlerFunc
	Options http.HandlerFunc
	Patch   http.HandlerFunc
	Head    http.HandlerFunc
}

func (rest *RESTEndPoint) Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && rest.Get != nil {
		rest.Get(w, r)
	} else if r.Method == http.MethodPost && rest.Post != nil {
		rest.Post(w, r)
	} else if r.Method == http.MethodPut && rest.Put != nil {
		rest.Put(w, r)
	} else if r.Method == http.MethodDelete && rest.Delete != nil {
		rest.Delete(w, r)
	} else if r.Method == http.MethodPatch && rest.Patch != nil {
		rest.Delete(w, r)
	} else if r.Method == http.MethodHead && rest.Head != nil {
		rest.Delete(w, r)
	} else if r.Method == http.MethodOptions && rest.Options != nil {
		rest.Options(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

package webjacker

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	EventParam        string = "_evt"
	SelfPath          string = "%s/%s"
	SelfPathForEvent  string = "%s/%s?" + EventParam + "=%s"
	UnsupportedEvent  string = "Unsupported event: %s"
	UnsupportedMethod string = "Unsupported method: %s"
)

type ResourceHandler func(w http.ResponseWriter, r *http.Request, params url.Values)

type HttpResource struct {
	Id       string
	pathBase string
	handlers map[string]ResourceHandler
}

func NewHttpResource(id string) *HttpResource {
	return &HttpResource{
		Id:       id,
		handlers: map[string]ResourceHandler{},
	}
}

// request.ParseForm() will already have been called when the handler is run
func (h *HttpResource) On(event string, handler ResourceHandler) *HttpResource {
	if h.handlers == nil {
		h.handlers = map[string]ResourceHandler{}
	}
	h.handlers[event] = handler
	return h
}

func (h *HttpResource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := GetParams(r)
	if evt := params.Get(EventParam); evt != "" {
		h.runHandler(evt, params, w, r, UnsupportedEvent)
	} else {
		h.runHandler(r.Method, params, w, r, UnsupportedMethod)
	}
}

func (h *HttpResource) runHandler(event string, params url.Values, w http.ResponseWriter, r *http.Request, errorString string) {
	if h.handlers == nil {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, errorString, event)
		return
	}

	if handler, ok := h.handlers[event]; ok {
		handler(w, r, params)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, errorString, event)
	}
}

func (h *HttpResource) Path() string {
	return fmt.Sprintf(SelfPath, h.pathBase, h.Id)
}

func (h *HttpResource) EventPath(event string) string {
	return fmt.Sprintf(SelfPathForEvent, h.pathBase, h.Id, event)
}

func RegisterHttpResource(resource *HttpResource) {
	RegisterHttpResourceOnPath(resource, http.DefaultServeMux, "")
}

func RegisterHttpResourceOnPath(resource *HttpResource, mux *http.ServeMux, pathBase string) {
	resource.pathBase = pathBase
	if pathBase != "" && !strings.HasPrefix(pathBase, "/") {
		resource.pathBase += "/"
	}

	mux.Handle(resource.Path(), resource)
}

func UnregisterHttpResource(resource *HttpResource, mux *http.ServeMux) {
	mux.Handle(resource.Path(), http.NotFoundHandler())
}

func GetParams(r *http.Request) url.Values {
	r.ParseForm()
	var values url.Values = r.Form
	for key, val := range r.Header {
		if strings.HasPrefix(key, "HX-") ||
			strings.HasPrefix(key, "@") {
			values.Add(key, val[0])
		}
	}
	return values
}

package app

import "net/http"

type Handler struct {
	*App
	Handler func(*Context)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := h.App.newContext(w, r)
	h.Handler(c)
}

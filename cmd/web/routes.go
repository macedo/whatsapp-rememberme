package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(RecoverPanic)

	// static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

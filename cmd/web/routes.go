package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/macedo/whatsapp-rememberme/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(RecoverPanic)

	mux.Get("/", handlers.HomePageHandler)

	// static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

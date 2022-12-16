package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/macedo/whatsapp-rememberme/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(SessionLoad)
	mux.Use(NoSurf)

	mux.Get("/", handlers.SignInPageHandler)
	mux.Post("/", handlers.SignInHandler)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.AdminPageHandler)
	})

	// static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

package handlers

import (
	"net/http"

	"github.com/macedo/whatsapp-rememberme/internal/app"
)

func App() *app.App {
	a := app.New()

	//a.Router().HandleFunc("/connect", ConnectHandler)

	a.Router().Handle("/", app.Handler{App: a, Handler: SignInPageHandler}).Methods("GET")
	a.Router().Handle("/", app.Handler{a, SignInHandler}).Methods("POST")
	//a.Router().Handle("/sign_out", app.Handler{a, SignOutHandler}).Methods("GET")

	// mux.Route("/admin", func(mux chi.Router) {
	// 	mux.Get("/", AdminPageHandler)

	// 	mux.Get("/devices/new", NewDevicePageHandler)
	// })

	fileServer := http.FileServer(http.Dir("./web/static"))
	a.Router().PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return a
}

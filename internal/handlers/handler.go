package handlers

import (
	"net/http"

	"github.com/macedo/whatsapp-rememberme/internal/config"
)

var app *config.AppConfig

func NewHandlers(a *config.AppConfig) {
	app = a
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Views.ExecuteTemplate(w, "home", "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

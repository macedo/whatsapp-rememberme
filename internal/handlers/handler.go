package handlers

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/macedo/whatsapp-rememberme/internal/config"
	"github.com/macedo/whatsapp-rememberme/internal/store"
)

var app *config.AppConfig
var repo store.Repository

func NewHandlers(r store.Repository, a *config.AppConfig) {
	app = a
	repo = r
}

type TemplateData struct {
	CSRFToken string
	Flash     string
}

func renderPage(w http.ResponseWriter, r *http.Request, tmplName, layoutName string) error {
	td := TemplateData{
		CSRFToken: nosurf.Token(r),
		Flash:     app.Session.PopString(r.Context(), "flash"),
	}

	return app.Views.ExecuteTemplate(w, tmplName, layoutName, td)
}

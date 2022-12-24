package handlers

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/kataras/blocks"
)

var views *blocks.Blocks

func init() {
	views = blocks.New("./web/views")
	if err := views.Load(); err != nil {
		log.Fatal(err)
	}
}

type TemplateData struct {
	CSRFToken string
	Flash     string
}

func renderPage(w http.ResponseWriter, r *http.Request, tmplName, layoutName string) error {
	td := TemplateData{
		CSRFToken: nosurf.Token(r),
		Flash:     session.PopString(r.Context(), "flash"),
	}

	return views.ExecuteTemplate(w, tmplName, layoutName, td)
}

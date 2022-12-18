package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/store"
	"github.com/macedo/whatsapp-rememberme/pkg/hash"
)

var encryptor hash.Encryptor

var repo store.Repository

var session *scs.SessionManager

var views *blocks.Blocks

func Init(r store.Repository, s *scs.SessionManager, v *blocks.Blocks, e hash.Encryptor) {
	encryptor = e
	repo = r
	session = s
	views = v
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

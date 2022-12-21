package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/domain/repository"
	"github.com/macedo/whatsapp-rememberme/pkg/hash"
	"go.mau.fi/whatsmeow"
	wasqlstore "go.mau.fi/whatsmeow/store/sqlstore"
)

var encryptor hash.Encryptor

var repo repository.DatabaseRepo

var session *scs.SessionManager

var views *blocks.Blocks

var waClients map[string]*whatsmeow.Client

var waContainer *wasqlstore.Container

func Init(
	r repository.DatabaseRepo,
	s *scs.SessionManager,
	v *blocks.Blocks,
	e hash.Encryptor,
	cli map[string]*whatsmeow.Client,
	c *wasqlstore.Container,
) {
	encryptor = e
	repo = r
	session = s
	views = v
	waClients = cli
	waContainer = c
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

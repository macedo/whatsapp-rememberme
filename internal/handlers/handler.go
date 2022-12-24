package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/macedo/whatsapp-rememberme/internal/domain/repository"
	"go.mau.fi/whatsmeow"
	wasqlstore "go.mau.fi/whatsmeow/store/sqlstore"
)

var repo repository.DatabaseRepo

var session *scs.SessionManager

var waClients map[string]*whatsmeow.Client

var waContainer *wasqlstore.Container

func Router() http.Handler {
	mux := chi.NewRouter()

	mux.HandleFunc("/connect", ConnectHandler)

	mux.Get("/", SignInPageHandler)
	mux.Post("/", SignInHandler)
	mux.Get("/sign_out", SignOutHandler)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Get("/", AdminPageHandler)

		mux.Get("/devices/new", NewDevicePageHandler)
	})

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

func Init(
	r repository.DatabaseRepo,
	s *scs.SessionManager,
	cli map[string]*whatsmeow.Client,
	c *wasqlstore.Container,
) {
	repo = r
	session = s
	waClients = cli
	waContainer = c
}

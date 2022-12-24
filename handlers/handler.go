package handlers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/macedo/whatsapp-rememberme/internal/domain/repository"
	"github.com/macedo/whatsapp-rememberme/pkg/hash"
	"go.mau.fi/whatsmeow"
	wasqlstore "go.mau.fi/whatsmeow/store/sqlstore"
)

var encryptor hash.Encryptor

var repo repository.DatabaseRepo

var session *scs.SessionManager

var waClients map[string]*whatsmeow.Client

var waContainer *wasqlstore.Container

func Init(
	r repository.DatabaseRepo,
	s *scs.SessionManager,
	e hash.Encryptor,
	cli map[string]*whatsmeow.Client,
	c *wasqlstore.Container,
) {
	encryptor = e
	repo = r
	session = s
	waClients = cli
	waContainer = c
}

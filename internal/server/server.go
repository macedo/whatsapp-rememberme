package server

import (
	"context"
	"database/sql"

	"github.com/macedo/whatsapp-rememberme/internal/store/sqlstore"
	"github.com/rs/zerolog/log"
)

type Server struct {
	db     *sql.DB
	name   string
	ctx    context.Context
	cancel context.CancelFunc
}

func New(db *sql.DB) *Server {
	return &Server{
		db:   db,
		name: "server",
	}
}

func (s *Server) Start() <-chan error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	errc := make(chan error)
	go func() {
		defer close(errc)
		if err := s.run(); err != nil {
			errc <- err
		}
	}()

	return errc
}

func (s *Server) Stop() {
	log.Info().Str("service", s.name).Msg("service stopped")
	s.cancel()
}

func (s *Server) run() error {
	log.Info().Str("service", s.name).Msg("service started")

	container := sqlstore.NewWithDB(s.db, "sqlite3", log.With().Str("component", "database").Logger())
	if err := container.Upgrade(); err != nil {
		return err
	}

	<-s.ctx.Done()
	return nil
}

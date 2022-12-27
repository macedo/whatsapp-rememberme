package app

import (
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type App struct {
	Options

	router *mux.Router
}

func (a *App) Router() *mux.Router {
	return a.router
}

func New() *App {
	opts := options_with_default(Options{})
	return &App{
		Options: opts,
		router:  mux.NewRouter(),
	}
}

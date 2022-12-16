package pqstore

import (
	"database/sql"

	"github.com/macedo/whatsapp-rememberme/internal/config"
	"github.com/macedo/whatsapp-rememberme/internal/driver"
)

var app *config.AppConfig

type Container struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewWithDB(db *driver.DB, a *config.AppConfig) *Container {
	app = a
	return &Container{
		App: a,
		DB:  db.SQL,
	}
}

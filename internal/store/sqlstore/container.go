package sqlstore

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Container struct {
	db      *sql.DB
	dialect string
	log     zerolog.Logger
}

func NewWithDB(db *sql.DB, dialect string, log zerolog.Logger) *Container {
	return &Container{
		db:      db,
		dialect: dialect,
		log:     log,
	}
}

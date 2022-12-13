package sqlstore

import (
	"database/sql"

	"github.com/macedo/whatsapp-rememberme/internal/logadapter"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Container struct {
	*sqlstore.Container
}

func NewWithDB(db *sql.DB, dialect string) *Container {
	dbLog := log.With().Str("module", "database").Logger()
	return &Container{sqlstore.NewWithDB(db, dialect, logadapter.WALogAdapter(dbLog))}
}

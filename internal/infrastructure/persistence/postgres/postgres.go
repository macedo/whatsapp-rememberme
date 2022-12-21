package postgres

import (
	"database/sql"

	"github.com/macedo/whatsapp-rememberme/internal/domain/repository"
)

type postgresRepo struct {
	db *sql.DB
}

func NewRepo(conn *sql.DB) repository.DatabaseRepo {
	return &postgresRepo{
		db: conn,
	}
}

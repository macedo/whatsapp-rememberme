package pqstore

import (
	"database/sql"

	"github.com/macedo/whatsapp-rememberme/pkg/hash"
)

type Container struct {
	DB        *sql.DB
	Encryptor hash.Encryptor
}

func NewWithDB(db *sql.DB, envcryptor hash.Encryptor) *Container {
	return &Container{
		DB:        db,
		Encryptor: envcryptor,
	}
}

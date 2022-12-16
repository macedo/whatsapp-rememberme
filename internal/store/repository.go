package store

import "github.com/macedo/whatsapp-rememberme/internal/models"

type Repository interface {
	// Users
	GetUserByUsername(username string) (models.User, error)

	Upgrade() error
}

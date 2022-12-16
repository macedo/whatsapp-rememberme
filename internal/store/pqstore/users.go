package pqstore

import (
	"context"
	"log"
	"time"

	"github.com/macedo/whatsapp-rememberme/internal/models"
)

const getUserByUsernameQuery = `
SELECT id, username, encrypted_password
FROM users
WHERE username = $1
`

func (c *Container) GetUserByUsername(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	row := c.DB.QueryRowContext(ctx, getUserByUsernameQuery, username)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.EncryptedPassword,
	)

	if err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}

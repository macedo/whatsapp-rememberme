package postgres

import (
	"context"
	"time"

	"github.com/macedo/whatsapp-rememberme/internal/domain/entity"
)

const getUserByUsernameQuery = `
SELECT id, username, encrypted_password
FROM users
WHERE username = $1
LIMIT 1
`

func (r *postgresRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, getUserByUsernameQuery, username)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.EncryptedPassword,
	)

	return user, err
}

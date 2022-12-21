package repository

import (
	"context"

	"github.com/macedo/whatsapp-rememberme/internal/domain/entity"
)

type DatabaseRepo interface {
	//public.users
	GetUserByUsername(context.Context, string) (entity.User, error)
}

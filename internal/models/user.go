package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID
	EncryptedPassword string
	Username          string
}

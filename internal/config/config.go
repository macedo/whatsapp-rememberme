package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/driver"
	"github.com/macedo/whatsapp-rememberme/internal/hash"
)

type AppConfig struct {
	DB           *driver.DB
	Domain       string
	Encryptor    hash.Encryptor
	IsProduction bool
	Session      *scs.SessionManager
	Views        *blocks.Blocks
}

func New() *AppConfig {
	return &AppConfig{}
}

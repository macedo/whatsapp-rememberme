package config

import (
	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/driver"
)

type AppConfig struct {
	DB    *driver.DB
	Views *blocks.Blocks
}

func New() *AppConfig {
	return &AppConfig{}
}

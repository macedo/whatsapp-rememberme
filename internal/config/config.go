package config

import (
	"github.com/macedo/whatsapp-rememberme/internal/driver"
)

type AppConfig struct {
	DB *driver.DB
}

func New() *AppConfig {
	return &AppConfig{}
}

package config

type AppConfig struct {
	DatabaseURL string
}

func New() *AppConfig {
	return &AppConfig{}
}

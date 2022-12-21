package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/macedo/whatsapp-rememberme/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfg *viper.Viper

var log *logrus.Logger

const version = "0.0.1"

func init() {
	gob.Register(uuid.UUID{})
	_ = os.Setenv("TZ", "America/Sao_Paulo")
}

func main() {
	log = logrus.New()

	// print info
	fmt.Printf("******************************************\n")
	fmt.Printf("** %sWhatsApp RememberMe%s v%sfmt built in %s\n", "\033[31m", "\033[0m", version, runtime.Version())
	fmt.Printf("**----------------------------------------\n")
	fmt.Printf("** Running with %d Processors\n", runtime.NumCPU())
	fmt.Printf("** Running on %s\n", runtime.GOOS)
	fmt.Printf("******************************************\n")

	setEnvironment()

	err := app.Run(cfg)
	if err != app.ErrShutdown {
		log.Fatalf("service stopped - %s", err)
	}

	log.Infof("service shutdown - %s", err)
}

func setEnvironment() {
	cfg = viper.New()
	cfg.AddConfigPath("./conf")
	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()

	parseOptions()
}

func parseOptions() {
	defaultOptions()

	err := cfg.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warnf("no config file found, loaded config from environment - default path ./config")
		default:
			log.Fatalf("error when fetching configuration - %s", err)
		}
	}
}

func defaultOptions() {
	cfg.SetDefault("app_env", "development")
	cfg.SetDefault("database_url", "postgres://postgres:postgres@localhost:15432/whatsapp_rememberme?sslmode=disable&timezone=UTC&connect_timeout=5")
	cfg.SetDefault("domain", "localhost")
	cfg.SetDefault("isProduction", false)
	cfg.SetDefault("listen_addr", ":8080")
	cfg.SetDefault("secret_key_base", "717c354e8cbfe8aba83e0897b2489314259d91805b02e9ca97e7a9e926406bd1119ec257cf3e77e7a2b8757e5041c49414523beb5086ccb406b14cbab5170535")
}

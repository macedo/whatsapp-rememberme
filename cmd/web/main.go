package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/macedo/whatsapp-rememberme/internal/handlers"
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/persistence"
	"github.com/macedo/whatsapp-rememberme/pkg/env"
)

const version = "0.0.1"

var ENV = env.Get("APP_ENV", "development")

func init() {
	_ = os.Setenv("TZ", "America/Sao_Paulo")
}

func main() {
	// print info
	fmt.Printf("******************************************\n")
	fmt.Printf("** %sWhatsApp RememberMe%s v%sfmt built in %s\n", "\033[31m", "\033[0m", version, runtime.Version())
	fmt.Printf("**----------------------------------------\n")
	fmt.Printf("** Running with %d Processors\n", runtime.NumCPU())
	fmt.Printf("** Running on %s\n", runtime.GOOS)
	fmt.Printf("******************************************\n")

	if err := persistence.LoadConfigFile(); err != nil {
		log.Fatal(err)
	}

	conn := persistence.Connections[ENV]
	if err := conn.Open(); err != nil {
		log.Fatal(err)
	}

	app := handlers.App()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

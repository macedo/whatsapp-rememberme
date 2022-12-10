package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/macedo/whatsapp-rememberme/internal/app"
)

func main() {
	err := app.Run()
	if err != nil && err != app.ErrShutdown {
		log.Fatalf("service stopped - %s", err)
	}
}

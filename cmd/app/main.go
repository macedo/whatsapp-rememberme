package main

import (
	"fmt"
	"os"

	"github.com/macedo/whatsapp-rememberme/internal/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "service stopped unexpectedly- %s", err)
		os.Exit(1)
	}

	fmt.Println("bye....")
}

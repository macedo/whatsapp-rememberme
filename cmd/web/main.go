package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/macedo/whatsapp-rememberme/internal/config"
)

var app config.AppConfig
var session *scs.SessionManager

const version = "0.0.1"

func init() {
	gob.Register(uuid.UUID{})
	_ = os.Setenv("TZ", "America/Sao_Paulo")
}

func main() {
	port, err := setupApp()
	if err != nil {
		log.Fatal(err)
	}

	// close channels & db when application ends
	defer app.DB.SQL.Close()

	// print info
	log.Printf("******************************************")
	log.Printf("** %sWhatsApp RememberMe%s v%s built in %s", "\033[31m", "\033[0m", version, runtime.Version())
	log.Printf("**----------------------------------------")
	log.Printf("** Running with %d Processors", runtime.NumCPU())
	log.Printf("** Running on %s", runtime.GOOS)
	log.Printf("******************************************")

	// create http server
	srv := &http.Server{
		Addr:              *port,
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Printf("starting HTTP server on port %s....", *port)

	// start server
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

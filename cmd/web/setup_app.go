package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/config"
	"github.com/macedo/whatsapp-rememberme/internal/driver"
	"github.com/macedo/whatsapp-rememberme/internal/handlers"
)

func setupApp() (*string, error) {
	port := flag.String("port", ":3000", "port to listen on")

	// database
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbUser := flag.String("dbuser", "postgres", "database user")
	dbPassword := flag.String("dbpassword", "postgres", "database password")
	dbPort := flag.Int("dbport", 5432, "database port")
	dbName := flag.String("dbname", "whatsapp_rememberme", "database name")
	dbSsl := flag.String("dbssl", "disable", "database ssl settings")

	flag.Parse()

	log.Println("connection to database....")
	dsnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
		*dbHost,
		*dbPort,
		*dbUser,
		*dbPassword,
		*dbName,
		*dbSsl,
	)

	db, err := driver.ConnectPostgres(dsnString)
	if err != nil {
		log.Fatal("cannot connect to database - ", err)
	}

	views := blocks.New("./views").Reload(true)
	if err := views.Load(); err != nil {
		log.Fatal("cannot load templates - ", err)
	}

	app = config.AppConfig{
		DB:    db,
		Views: views,
	}

	handlers.NewHandlers(&app)

	return port, nil
}

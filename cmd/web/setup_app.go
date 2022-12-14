package main

import (
	"flag"

	"github.com/macedo/whatsapp-rememberme/internal/config"
)

func setupApp() (*string, error) {
	port := flag.String("port", ":3000", "port to listen on")
	//databaseURL := flag.String("database-url", "", "postgres connection url")

	flag.Parse()

	// log.Println("connection to database....")
	// db, err := driver.ConnectPostgres(*databaseURL)
	// if err != nil {
	// 	log.Fatal("cannot connect to database", err)
	// }

	a := config.AppConfig{
		//DB: db,
	}

	app = a

	return port, nil
}

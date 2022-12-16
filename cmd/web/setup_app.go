package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/config"
	"github.com/macedo/whatsapp-rememberme/internal/driver"
	"github.com/macedo/whatsapp-rememberme/internal/handlers"
	"github.com/macedo/whatsapp-rememberme/internal/hash"
	"github.com/macedo/whatsapp-rememberme/internal/store/pqstore"
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
	domain := flag.String("domain", "localhost", "domain name (e.g. example.com)")
	inProduction := flag.Bool("production", false, "application is in production")

	flag.Parse()

	encryptor := hash.NewEncryptor("pepper")

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

	log.Println("initializing template engine")

	views := blocks.New("./views").Reload(true)
	if err := views.Load(); err != nil {
		log.Fatal("cannot load templates - ", err)
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour

	app = config.AppConfig{
		DB:           db,
		Domain:       *domain,
		Encryptor:    encryptor,
		IsProduction: *inProduction,
		Session:      session,
		Views:        views,
	}

	log.Println("initializing store")
	container := pqstore.NewWithDB(db, &app)
	if err := container.Upgrade(); err != nil {
		log.Fatal("failed to upgrade database - ", err)
	}

	handlers.NewHandlers(container, &app)

	return port, nil
}

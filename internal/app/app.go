package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/macedo/whatsapp-rememberme/internal/handlers"
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/persistence"
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/persistence/postgres"
	"github.com/procyon-projects/chrono"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mau.fi/whatsmeow"
	wasqlstore "go.mau.fi/whatsmeow/store/sqlstore"
)

var (
	ErrShutdown = fmt.Errorf("application shutdown gracefully")
)

var cfg *viper.Viper

var db *sql.DB

var log *logrus.Logger

var runCancel context.CancelFunc

var runCtx context.Context

var session *scs.SessionManager

var scheduler chrono.TaskScheduler

var srv *http.Server

var waClients map[string]*whatsmeow.Client

var waContainer *wasqlstore.Container

func Run(c *viper.Viper) (err error) {
	runCtx, runCancel = context.WithCancel(context.Background())

	cfg = c

	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)

	secret_key_base := c.GetString("secret_key_base")
	if secret_key_base == "" {
		return fmt.Errorf("secret_key_base need to be informed")
	}
	scheduler = chrono.NewDefaultTaskScheduler()
	defer scheduler.Shutdown()

	log.Infof("connecting to postgresql")
	err = persistence.LoadConfigFile()
	if err != nil {
		return err
	}
	conn := persistence.Connections[c.GetString("app_env")]
	fmt.Printf("%v", conn)
	db, err = conn.Open()
	if err != nil {
		return fmt.Errorf("could not establish postgres connection - %s", err)
	}
	defer db.Close()

	//dbrepo
	_ = postgres.NewRepo(db)

	log.Infof("initializing session manager with postgresql backend")
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(db)
	session.Cookie.SameSite = http.SameSiteLaxMode

	waContainer = wasqlstore.NewWithDB(db, "postgres", nil)
	if err = waContainer.Upgrade(); err != nil {
		return fmt.Errorf("could not upgrade whatsmeow database to last version - %s", err)
	}

	devices, err := waContainer.GetAllDevices()
	if err != nil {
		return fmt.Errorf("could not get device - %s", err)
	}
	log.Debugf("found %d device(s)", len(devices))

	waClients = make(map[string]*whatsmeow.Client, len(devices))
	for _, device := range devices {
		waClients[device.ID.String()] = whatsmeow.NewClient(device, nil)
		defer waClients[device.ID.String()].Disconnect()
	}

	srv = &http.Server{
		Addr:              cfg.GetString("listen_addr"),
		Handler:           handlers.Router(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	// Kick off gracefull shutdown go routine
	go func() {
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGTERM, syscall.SIGINT)

		s := <-trap
		log.Infof("received shutdown signal - %s", s)
		Stop()
	}()

	log.Infof("starting HTTP listener on %s", cfg.GetString("listen_addr"))
	if err = srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			<-runCtx.Done()
			return ErrShutdown
		}

		return fmt.Errorf("unable to start HTTP server - %s", err)
	}

	return nil
}

func Stop() {
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Errorf("unexpected errow while shutting down HTTP server - %s", err)
	}
	defer runCancel()
}

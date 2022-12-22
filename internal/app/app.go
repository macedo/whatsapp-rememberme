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
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kataras/blocks"
	"github.com/macedo/whatsapp-rememberme/internal/handlers"
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/persistence"
	"github.com/macedo/whatsapp-rememberme/internal/infrastructure/persistence/postgres"
	"github.com/macedo/whatsapp-rememberme/pkg/hash"
	"github.com/macedo/whatsapp-rememberme/pkg/middleware"
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

var encryptor hash.Encryptor

var log *logrus.Logger

var runCancel context.CancelFunc

var runCtx context.Context

var session *scs.SessionManager

var scheduler chrono.TaskScheduler

var srv *http.Server

var views *blocks.Blocks

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
	encryptor = hash.NewEncryptor(secret_key_base)

	scheduler = chrono.NewDefaultTaskScheduler()
	defer scheduler.Shutdown()

	log.Infof("initializing template engine")
	views = blocks.New("./web/views").Reload(true)
	if err = views.Load(); err != nil {
		return fmt.Errorf("could not load views - %s", err)
	}

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

	dbrepo := postgres.NewRepo(db)

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

	handlers.Init(dbrepo, session, views, encryptor, waClients, waContainer)

	router := initialize_router()

	srv = &http.Server{
		Addr:              cfg.GetString("listen_addr"),
		Handler:           router,
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

func initialize_router() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(SessionLoad)
	mux.Use(NoSurf)

	mux.HandleFunc("/connect", handlers.ConnectHandler)

	mux.Get("/", handlers.SignInPageHandler)
	mux.Post("/", handlers.SignInHandler)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.AdminPageHandler)

		mux.Get("/devices/new", handlers.NewDevicePageHandler)
	})

	// static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

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
	"github.com/macedo/whatsapp-rememberme/internal/store/pqstore"
	"github.com/macedo/whatsapp-rememberme/pkg/hash"
	"github.com/macedo/whatsapp-rememberme/pkg/middleware"
	"github.com/procyon-projects/chrono"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func Run(c *viper.Viper) (err error) {
	runCtx, runCancel = context.WithCancel(context.Background())

	cfg = c

	log = logrus.New()

	secret_key_base := c.GetString("secret_key_base")
	if secret_key_base == "" {
		return fmt.Errorf("secret_key_base need to be informed")
	}
	encryptor = hash.NewEncryptor(secret_key_base)

	scheduler = chrono.NewDefaultTaskScheduler()
	defer scheduler.Shutdown()

	log.Infof("initializing template engine")
	views = blocks.New("./views").Reload(true)
	if err = views.Load(); err != nil {
		return fmt.Errorf("could not load views - %s", err)
	}

	log.Infof("connecting to postgresql")
	db, err = sql.Open("pgx", cfg.GetString("database_url"))
	if err != nil {
		return fmt.Errorf("could not establish postgres db connection - %s", err)
	}

	log.Infof("initializing session manager with postgresql backend")
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(db)

	log.Infof("initializing store")
	container := pqstore.NewWithDB(db, encryptor)
	if err := container.Upgrade(); err != nil {
		return fmt.Errorf("could not upgrade database to last version - %s", err)
	}

	handlers.Init(container, session, views, encryptor)

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
		signal.Notify(trap, syscall.SIGTERM)

		s := <-trap
		log.Infof("received shutdown signal %s", s)
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
	mux.Use(middleware.RequestLogger(log))
	mux.Use(SessionLoad)
	mux.Use(NoSurf)

	mux.Get("/", handlers.SignInPageHandler)
	mux.Post("/", handlers.SignInHandler)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.AdminPageHandler)
	})

	// static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

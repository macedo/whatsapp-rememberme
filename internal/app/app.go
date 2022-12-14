package app

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/macedo/whatsapp-rememberme/internal/handler"
	"github.com/macedo/whatsapp-rememberme/internal/server"
	"github.com/macedo/whatsapp-rememberme/internal/whatsapp"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/br"
	"github.com/olebedev/when/rules/common"
	"github.com/procyon-projects/chrono"
	"github.com/rs/zerolog/log"
)

func Run() error {
	// Setup signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	dsn := "file:wpp_store.db?_foreign_keys=on"
	log.Info().Str("dsn", dsn).Msg("connecting to database")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	log.Info().Msg("starting task scheduler")
	scheduler := chrono.NewDefaultTaskScheduler()
	defer scheduler.Shutdown()

	log.Info().Msg("configuring parser")
	w := when.New(&rules.Options{
		Distance: 10,
	})
	w.Add(common.All...)
	w.Add(br.All...)

	evtHandler := handler.NewEventHandler(w, scheduler)

	wa := whatsapp.New(db, evtHandler)
	defer wa.Stop()
	waCh := wa.Start()

	srv := server.New(db)
	defer srv.Stop()
	srvCh := srv.Start()

	select {
	case err := <-srvCh:
		log.Error().Err(err).Msg("server error")
	case err := <-waCh:
		log.Error().Err(err).Msg("whatstapp client error")
	case sig := <-sigCh:
		log.Info().Msgf("got %v signal", sig)
		break
	}

	log.Info().Msg("shutting down")

	return nil
}

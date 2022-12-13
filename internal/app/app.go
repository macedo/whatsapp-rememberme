package app

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/macedo/whatsapp-rememberme/internal/handler"
	"github.com/macedo/whatsapp-rememberme/internal/whatsapp"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/br"
	"github.com/olebedev/when/rules/common"
	"github.com/procyon-projects/chrono"
)

func Run() error {
	// Setup signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("connecting to database...")
	db, err := sql.Open("sqlite3", "file:wpp_store.db?_foreign_keys=on")
	if err != nil {
		return err
	}

	log.Printf("starting task scheduler...")
	scheduler := chrono.NewDefaultTaskScheduler()
	defer scheduler.Shutdown()

	log.Printf("configuring parser...")
	w := when.New(&rules.Options{
		Distance: 10,
	})
	w.Add(common.All...)
	w.Add(br.All...)

	evtHandler := handler.NewEventHandler(w, scheduler)

	wa := whatsapp.New(db, evtHandler)
	waCh := wa.Start()

	select {
	case err := <-waCh:
		log.Printf("error: %v", err)
	case sig := <-sigCh:
		log.Printf("got %v signal", sig)
		break
	}

	log.Printf("shutting down")
	wa.Stop()

	return nil
}

package app

import (
	"fmt"
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

// Common errors returned by this app.
var (
	ErrShutdown = fmt.Errorf("application was shutdown gracefully")
)

func Run() error {
	// Setup signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	scheduler := chrono.NewDefaultTaskScheduler()
	defer scheduler.Shutdown()

	w := when.New(&rules.Options{
		Distance: 10,
	})
	w.Add(common.All...)
	w.Add(br.All...)

	evtHandler := handler.NewEventHandler(w, scheduler)

	wa := whatsapp.New(evtHandler)
	waCh := wa.Start()

	select {
	case err := <-waCh:
		log.Printf("error: %v", err)
		break
	case sig := <-sigCh:
		log.Printf("got signal %v", sig)
		break
	}

	log.Printf("shutting down")
	wa.Stop()

	return <-waCh
}

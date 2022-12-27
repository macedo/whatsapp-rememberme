package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func (a *App) Run() error {
	var wg sync.WaitGroup

	srv := http.Server{
		Addr:              a.Addr,
		Handler:           a.router,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(a.Context, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		ctx, cfn := context.WithTimeout(context.Background(), time.Duration(60*time.Second))
		defer cfn()

		//"shutting down application"
		if err := srv.Shutdown(ctx); err != nil {
			//"shutting down server: ", er
		}
		cfn()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.Stop(srv.ListenAndServe())
	}()

	wg.Wait()
	// "shutdown complete"
	err := ctx.Err()
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}

func (a *App) Stop(err error) error {
	ce := a.Context.Err()
	if ce != nil {
		return errors.New("application has already been canceled")
	}

	a.cancel()

	return nil
}

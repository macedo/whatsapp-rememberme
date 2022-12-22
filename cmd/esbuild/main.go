package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()
		stop()
	}()

	run(ctx)
}

func run(ctx context.Context) {
	result := api.Build(api.BuildOptions{
		Color:    api.ColorAlways,
		LogLevel: api.LogLevelInfo,
		EntryPoints: []string{
			"web/src/admin.js",
		},
		EntryNames: "[dir]/[name].bundle",
		Outdir:     "web/static",
		Bundle:     true,
		Write:      true,
		Watch:      &api.WatchMode{},
	})

	<-ctx.Done()
	result.Stop()
	fmt.Println("stopped watching")
}

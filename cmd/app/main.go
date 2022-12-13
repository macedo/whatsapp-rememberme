package main

import (
	"flag"
	"os"

	"github.com/macedo/whatsapp-rememberme/internal/app"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var debug *bool

func init() {
	debug = flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := app.Run(); err != nil {
		log.Logger.Error().Err(err).Msg("service stopped unexpectedly")
		os.Exit(1)
	}

	log.Info().Msg("bye....")
}

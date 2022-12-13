package logadapter

import (
	"github.com/rs/zerolog"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type logAdapter struct {
	zerolog.Logger
}

func WALogAdapter(log zerolog.Logger) waLog.Logger {
	return &logAdapter{log}
}

func (l *logAdapter) Debugf(msg string, args ...interface{}) {
	l.Debug().Msgf(msg, args)
}

func (l *logAdapter) Errorf(msg string, args ...interface{}) {
	l.Error().Msgf(msg, args)
}

func (l *logAdapter) Infof(msg string, args ...interface{}) {
	l.Info().Msgf(msg, args)
}

func (l *logAdapter) Sub(module string) waLog.Logger {
	return WALogAdapter(l.With().Str("module", module).Logger())
}

func (l *logAdapter) Warnf(msg string, args ...interface{}) {
	l.Warn().Msgf(msg, args)
}

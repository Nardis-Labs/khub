package logging

import (
	"github.com/rs/zerolog"
)

var logLevels = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"trace": zerolog.TraceLevel,
	"error": zerolog.ErrorLevel,
}

// InitLogger sets up logging level, and log formatting
func InitLogger(logLevel string) {
	zerolog.SetGlobalLevel(logLevels[logLevel])
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

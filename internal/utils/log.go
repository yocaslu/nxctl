package utils

import (
	"log/slog"
	"os"
)

// global logger
var glogger *slog.Logger

func InitLogger(debug bool) {
	var level slog.Level = slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	glogger = slog.New(handler)

	slog.SetDefault(glogger)
}

// LogForModule cria um sub-logger isolado contendo o atributo do módulo correspondente
func LogForModule(moduleName string) *slog.Logger {
	if glogger == nil {
		InitLogger(false)
	}

	return glogger.With("module", moduleName)
}

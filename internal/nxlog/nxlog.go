package nxlog

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger(debug bool) {
	var level slog.Level = slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.New(handler)
	slog.SetDefault(logger)
}

// LogForModule cria um sub-logger isolado contendo o atributo do módulo correspondente
func ForModule(moduleName string) *slog.Logger {
	if logger == nil {
		InitLogger(false)
	}

	return logger.With("module", moduleName)
}

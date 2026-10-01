package nxlog

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(debug bool) {
	var level slog.Level = slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// LogForModule cria um sub-logger isolado contendo o atributo do módulo correspondente
func ForModule(moduleName string) *slog.Logger {
	if Logger == nil {
		InitLogger(false)
	}

	return Logger.With("module", moduleName)
}

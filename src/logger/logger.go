package logger

import (
	"log/slog"
	"os"
)

func SetUp(logLevel, logFormat string) *slog.Logger {
	switch logLevel {
	case "debug":
		switch logFormat {
		case "text":
			return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case "json":
			return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		default:
			panic("incorrect logFormat")
		}

	case "prod":
		switch logFormat {
		case "text":
			return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		case "json":
			return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		default:
			panic("incorrect logFormat")
		}
	default:
		panic("incorrect loglevel")
	}
}

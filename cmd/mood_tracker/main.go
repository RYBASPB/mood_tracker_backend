package main

import (
	"log/slog"
	"mood_tracker/internal/config"
	"mood_tracker/internal/storage/postgresql"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting logger", slog.String("env", cfg.Env))
	log.Debug("log level is Debug")

	_, err := postgresql.ConnectToDB()
	if err != nil {
		log.Error("Error in initialization of DB", slog.AnyValue(err))
	}
	log.Info("connection to database is successful")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

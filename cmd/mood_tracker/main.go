package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"mood_tracker/internal/config"
	"mood_tracker/internal/http-server/handlers/score/add"
	httpLogger "mood_tracker/internal/http-server/middleware/logger"
	"mood_tracker/internal/storage/postgresql"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	db, err := postgresql.ConnectToDB()
	if err != nil {
		log.Error("Error in initialization of DB", slog.AnyValue(err))
	}
	log.Info("connection to database is successful")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(httpLogger.New(log))

	router.Route("/scores", func(r chi.Router) {
		r.Post("/", add.New(log, db))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("error in starting server", slog.AnyValue(err))
		}
	}()
	log.Info("server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ContextTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("error in shutting down", slog.AnyValue(err))
		return
	}
	log.Info("server stopped")

	db.CloseConnection()
	log.Info("db connection closed")
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

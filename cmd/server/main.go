package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"ipquery/internal/handlers"
	"ipquery/internal/ipinfo"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	svc := ipinfo.NewService()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	h := &handlers.IPHandler{
		Service: svc,
		Logger:  logger,
	}

	r := handlers.NewRouter(h)

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.SetDefault(logger)
	slog.Info("iniciando servidor", "port", os.Getenv("PORT"), "env", "production")

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("error crítico al iniciar servidor", "error", err)
		os.Exit(1)
	}
}

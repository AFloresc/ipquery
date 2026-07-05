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
	h := &handlers.IPHandler{Service: svc}

	r := handlers.NewRouter(h)

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Puedes cambiar a LevelDebug para desarrollo
	}))

	slog.SetDefault(logger)
	slog.Info("iniciando servidor", "port", os.Getenv("PORT"), "env", "production")

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("error crítico al iniciar servidor", "error", err)
		os.Exit(1)
	}
}

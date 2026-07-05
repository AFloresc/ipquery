package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	mylogger "ipquery/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

// NewRouter configura y retorna el enrutador
func NewRouter(h *IPHandler) http.Handler {

	r := chi.NewRouter()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(mylogger.NewStructuredLogger(logger))

	// Rutas
	r.Route("/v1", func(r chi.Router) {
		r.Get("/ip/{ip}", h.GetIPInfo)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	return r
}

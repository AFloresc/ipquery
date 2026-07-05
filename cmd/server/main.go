package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ipquery/internal/handlers"
	"ipquery/internal/ipinfo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()
	svc := ipinfo.NewService()
	h := &handlers.IPHandler{Service: svc}

	r := chi.NewRouter()

	// 1. Configuración de CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	// 2. Aplicamos CORS como primer middleware
	r.Use(c.Handler)

	// 3. Middlewares existentes
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// 4. Rutas
	r.Route("/v1", func(r chi.Router) {
		r.Get("/ip/{ip}", h.GetIPInfo)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("🚀 IP Service corriendo en http://localhost:" + port + " con CORS habilitado")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error crítico: %v", err)
	}
}

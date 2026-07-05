package main

import (
	"log"
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

	log.Printf("🚀 IP Service escuchando en el puerto %s con CORS habilitado", os.Getenv("PORT"))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error crítico: %v", err)
	}
}

package main

import (
	"ipquery/internal/handlers"
	"ipquery/internal/ipinfo"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// 1. Inicialización de dependencias
	// Creamos el servicio que maneja la lógica de API + Caché
	svc := ipinfo.NewService()

	// Inyectamos el servicio en el handler
	h := &handlers.IPHandler{Service: svc}

	// 2. Configuración del Router
	r := chi.NewRouter()

	// Middlewares esenciales para un servicio robusto
	r.Use(middleware.RequestID)                 // Asigna un ID único a cada petición (útil para logs)
	r.Use(middleware.RealIP)                    // Obtiene la IP real si estamos detrás de un proxy/load balancer
	r.Use(middleware.Logger)                    // Loguea peticiones en consola
	r.Use(middleware.Recoverer)                 // Evita que el servidor caiga ante errores fatales
	r.Use(middleware.Timeout(60 * time.Second)) // Timeout global para evitar bloqueos

	// 3. Rutas
	r.Route("/v1", func(r chi.Router) {
		r.Get("/ip/{ip}", h.GetIPInfo)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	})

	// 4. Obtener puerto de la variable de entorno (Render lo inyecta automáticamente)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback local
	}

	// 5. Configuración del servidor (Senior Approach)
	// Definimos explícitamente timeouts para evitar ataques de Slowloris
	srv := &http.Server{
		Addr:         ":" + port, // Usamos el puerto dinámico
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("🚀 IP Service corriendo en http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error crítico al iniciar servidor: %v", err)
	}
}

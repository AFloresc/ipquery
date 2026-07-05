package handlers

import (
	"encoding/json"
	"ipquery/internal/ipinfo"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type IPHandler struct {
	Service *ipinfo.Service
	Logger  *slog.Logger
}

func (h *IPHandler) GetIPInfo(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	// Delegamos la lógica al servicio
	data, err := h.Service.GetInfo(ip)
	if err != nil {
		h.Logger.Error("error al procesar ip",
			"ip", ip,
			"error", err,
			"request_id", r.Context().Value(middleware.RequestIDKey),
		)
		http.Error(w, "No se pudo procesar la IP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

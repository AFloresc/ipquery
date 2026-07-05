package handlers

import (
	"encoding/json"
	"ipquery/internal/ipinfo"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IPHandler struct {
	Service *ipinfo.Service
}

func (h *IPHandler) GetIPInfo(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	// Delegamos la lógica al servicio
	data, err := h.Service.GetInfo(ip)
	if err != nil {
		http.Error(w, "No se pudo procesar la IP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

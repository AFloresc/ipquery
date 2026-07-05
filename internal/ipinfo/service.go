package ipinfo

import (
	"sync"
	"time"
)

// IPData define la respuesta enriquecida
type IPData struct {
	IP      string `json:"ip"`
	ISP     string `json:"isp"`
	Country string `json:"country"`
	City    string `json:"city"`
	IsVPN   bool   `json:"is_vpn"`
	Risk    int    `json:"risk_score"`
}

// CacheEntry ayuda a expirar datos antiguos
type CacheEntry struct {
	Data      *IPData
	ExpiresAt time.Time
}

type Service struct {
	cache sync.Map // Caché en memoria para evitar llamadas redundantes
}

func NewService() *Service {
	return &Service{}
}

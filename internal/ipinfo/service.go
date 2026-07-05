package ipinfo

import (
	"sync"
	"time"
)

// --- Modelos unificados ---

type ISPInfo struct {
	ASN string `json:"asn"`
	Org string `json:"org"`
	ISP string `json:"isp"`
}

type LocationInfo struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Zipcode     string  `json:"zipcode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	Localtime   string  `json:"localtime"`
}

type RiskInfo struct {
	IsMobile     bool `json:"is_mobile"`
	IsVPN        bool `json:"is_vpn"`
	IsTor        bool `json:"is_tor"`
	IsProxy      bool `json:"is_proxy"`
	IsDatacenter bool `json:"is_datacenter"`
	RiskScore    int  `json:"risk_score"`
}

type IPData struct {
	IP       string       `json:"ip"`
	ISP      ISPInfo      `json:"isp"`
	Location LocationInfo `json:"location"`
	Risk     RiskInfo     `json:"risk"`
}

// --- Lógica de Servicio y Caché ---

type CacheEntry struct {
	Data      *IPData
	ExpiresAt time.Time
}

type Service struct {
	cache sync.Map
}

func NewService() *Service {
	return &Service{}
}

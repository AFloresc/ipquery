package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type APIResponse struct {
	IP       string `json:"ip"`
	Location struct {
		City      string  `json:"city"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"time_zone"` // Corregido: ip.guide usa "time_zone"
	} `json:"location"`
	Network struct {
		AS struct {
			ASN          int    `json:"asn"`
			Organization string `json:"organization"`
			Name         string `json:"name"`
		} `json:"autonomous_system"`
	} `json:"network"`
}

func (s *Service) GetInfo(ip string) (*IPData, error) {
	// 1. Verificación de caché
	if val, ok := s.cache.Load(ip); ok {
		entry := val.(CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data, nil
		}
	}

	// 2. Consulta a API
	url := fmt.Sprintf("https://ip.guide/%s", ip)
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "IPQuery-Service/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider returned status: %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// 3. Mapeo a IPData
	data := &IPData{
		IP: apiResp.IP,
		ISP: ISPInfo{
			ASN: fmt.Sprintf("AS%d", apiResp.Network.AS.ASN),
			Org: apiResp.Network.AS.Organization,
			ISP: apiResp.Network.AS.Name,
		},
		Location: LocationInfo{
			Country:     apiResp.Location.Country,
			CountryCode: "N/A", // Valor por defecto amigable
			City:        apiResp.Location.City,
			State:       "N/A",
			Zipcode:     "N/A",
			Latitude:    apiResp.Location.Latitude,
			Longitude:   apiResp.Location.Longitude,
			Timezone:    apiResp.Location.Timezone,
			Localtime:   time.Now().Format(time.RFC3339), // Formato estándar ISO
		},
		Risk: RiskInfo{
			IsMobile:     false,
			IsVPN:        false,
			IsTor:        false,
			IsProxy:      false,
			IsDatacenter: false,
			RiskScore:    0,
		},
	}

	// 4. Caché
	s.cache.Store(ip, CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})

	return data, nil
}

func detectRisk(org string, asn int) RiskInfo {
	risk := RiskInfo{}

	// Lógica simple de detección (puedes ampliar esta lista)
	if strings.Contains(org, "Cloud") || strings.Contains(org, "Amazon") {
		risk.IsDatacenter = true
		risk.RiskScore = 50
	}

	if strings.Contains(org, "Vodafone") || strings.Contains(org, "Movistar") {
		risk.IsMobile = true
	}

	return risk
}

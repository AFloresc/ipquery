package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// APIResponse ajustada para la estructura que devuelve https://ip.guide/
type APIResponse struct {
	Network struct {
		AutonomousSystemNumber       int    `json:"autonomous_system_number"`
		AutonomousSystemOrganization string `json:"autonomous_system_organization"`
	} `json:"network"`
	Location struct {
		Country struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"country"`
		City struct {
			Name string `json:"name"`
		} `json:"city"`
		TimeZone  string  `json:"time_zone"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Postal struct {
		Code string `json:"code"`
	} `json:"postal"`
}

func (s *Service) GetInfo(ip string) (*IPData, error) {
	// 1. Verificación de caché
	if val, ok := s.cache.Load(ip); ok {
		entry := val.(CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data, nil
		}
	}

	// 2. Consulta a API con cliente configurado
	url := fmt.Sprintf("https://ip.guide/%s", ip)
	client := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "IPQuery-Service/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// 3. Mapeo (Transformación de APIResponse de ip.guide a tu IPData)
	data := &IPData{
		IP: ip,
		ISP: ISPInfo{
			ASN: fmt.Sprintf("AS%d", apiResp.Network.AutonomousSystemNumber),
			Org: apiResp.Network.AutonomousSystemOrganization,
			ISP: apiResp.Network.AutonomousSystemOrganization,
		},
		Location: LocationInfo{
			Country:     apiResp.Location.Country.Name,
			CountryCode: apiResp.Location.Country.Code,
			City:        apiResp.Location.City.Name,
			State:       "", // ip.guide a veces no devuelve region/state
			Zipcode:     apiResp.Postal.Code,
			Latitude:    apiResp.Location.Latitude,
			Longitude:   apiResp.Location.Longitude,
			Timezone:    apiResp.Location.TimeZone,
			Localtime:   time.Now().Format("2006-01-02T15:04:05"),
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

	// 4. Guardar en caché
	s.cache.Store(ip, CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})

	return data, nil
}

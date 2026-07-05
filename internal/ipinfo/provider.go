package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// APIResponse mapea la respuesta gratuita de ip-api.com
type APIResponse struct {
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	RegionName  string  `json:"regionName"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	AS          string  `json:"as"`
}

func (s *Service) GetInfo(ip string) (*IPData, error) {
	// 1. Caché
	if val, ok := s.cache.Load(ip); ok {
		entry := val.(CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data, nil
		}
	}

	// 2. Consulta a API (ip-api.com no requiere token)
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,query,country,countryCode,city,regionName,zip,lat,lon,timezone,isp,as", ip)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status == "fail" {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	// 3. Lógica de Detección de Riesgo (Basada en ISP/ASN)
	risk := detectRisk(apiResp.ISP, apiResp.AS)

	// 4. Mapeo a IPData
	data := &IPData{
		IP: apiResp.Query,
		ISP: ISPInfo{
			ASN: apiResp.AS,
			ISP: apiResp.ISP,
			Org: apiResp.ISP,
		},
		Location: LocationInfo{
			Country:     apiResp.Country,
			CountryCode: apiResp.CountryCode,
			City:        apiResp.City,
			State:       apiResp.RegionName,
			Zipcode:     apiResp.Zip,
			Latitude:    apiResp.Lat,
			Longitude:   apiResp.Lon,
			Timezone:    apiResp.Timezone,
			Localtime:   time.Now().Format("2006-01-02T15:04:05"),
		},
		Risk: risk,
	}

	s.cache.Store(ip, CacheEntry{Data: data, ExpiresAt: time.Now().Add(1 * time.Hour)})
	return data, nil
}

// detectRisk analiza el nombre del proveedor para marcar riesgos
func detectRisk(isp, asn string) RiskInfo {
	info := strings.ToLower(isp + " " + asn)
	risk := RiskInfo{
		IsMobile:     false,
		IsVPN:        false,
		IsTor:        false,
		IsProxy:      false,
		IsDatacenter: false,
		RiskScore:    0,
	}

	// Lista de palabras clave detectables
	if strings.Contains(info, "vpn") || strings.Contains(info, "proxy") {
		risk.IsVPN = true
		risk.IsProxy = true
		risk.RiskScore = 70
	}
	if strings.Contains(info, "amazon") || strings.Contains(info, "google cloud") || strings.Contains(info, "digitalocean") || strings.Contains(info, "ovh") {
		risk.IsDatacenter = true
		risk.RiskScore = 50
	}
	if strings.Contains(info, "vodafone") || strings.Contains(info, "movistar") || strings.Contains(info, "orange") || strings.Contains(info, "t-mobile") {
		risk.IsMobile = true
	}

	return risk
}

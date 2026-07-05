package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Service) GetInfo(ip string) (*IPData, error) {
	// 1. Verificación de caché
	if val, ok := s.cache.Load(ip); ok {
		entry := val.(CacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Data, nil
		}
	}

	// 2. Consulta a API (ejemplo: ipapi.co)
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var data IPData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// 3. Guardar en caché por 1 hora
	s.cache.Store(ip, CacheEntry{
		Data:      &data,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})

	return &data, nil
}

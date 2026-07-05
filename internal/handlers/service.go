package ipinfo

type IPData struct {
	IP      string `json:"ip"`
	ISP     string `json:"isp"`
	City    string `json:"city"`
	Country string `json:"country"`
	// ... resto de campos
}

// Service define el contrato de nuestra lógica de negocio
type Service interface {
	GetInfo(ip string) (*IPData, error)
}

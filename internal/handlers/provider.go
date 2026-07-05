package ipinfo

import (
	"embed"

	"github.com/oschwald/geoip2-golang"
)

//go:embed ../../data/GeoLite2-City.mmdb
var dbFile embed.FS

type geoProvider struct {
	db *geoip2.Reader
}

func NewProvider() (*geoProvider, error) {
	data, err := dbFile.ReadFile("data/GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	db, err := geoip2.FromBytes(data)
	return &geoProvider{db: db}, err
}

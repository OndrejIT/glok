package geoip

import (
	log "github.com/Sirupsen/logrus"
	"net"
)

type LookupResponse struct {
	Country   string  `json:"coutry"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func Lookup(ip net.IP) LookupResponse {
	record, err := db.City(ip)
	if err != nil {
		log.Error(err)
	}
	return LookupResponse{
		Country:   record.Country.IsoCode,
		City:      record.City.Names["en"],
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}

}

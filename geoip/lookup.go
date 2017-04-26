package geoip

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"errors"
	"fmt"
	"net/http"
)

type LookupResponse struct {
	Country   string  `json:"coutry"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func Lookup(ip net.IP) (LookupResponse, int, error) {
	record, err := db.City(ip)
	if err != nil {
		log.Errorf("[Lookup] %s", err)
		return LookupResponse{}, http.StatusBadRequest, err
	}

	lookup := LookupResponse{
		Country:   record.Country.IsoCode,
		City:      record.City.Names["en"],
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}

	if lookup.Country == "" {
		err := errors.New(fmt.Sprintf("Lookup for %s not found.", ip))
		log.Errorf("[Lookup] %s", err)
		return LookupResponse{}, http.StatusNotFound, err
	}

	return lookup, http.StatusOK, nil

}

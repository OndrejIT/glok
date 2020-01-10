package geoip

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type LookupResponse struct {
	Country   string  `json:"country"`
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

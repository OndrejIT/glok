package api

import (
	"encoding/json"
	"fmt"
	conf "github.com/spf13/viper"
	"github.com/ondrejit/glok/geoip"
	"net"
	"net/http"
	"strings"
)

func V1(w http.ResponseWriter, r *http.Request) {
	ipString := r.FormValue("ip")
	if ipString == "" {
		ipString, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	lookup, status, err := geoip.Lookup(net.ParseIP(ipString))
	if err != nil {
		w.WriteHeader(status)
		return

	}

	switch r.URL.Path {
	case "/v1/lookup":
		writeLookup(w, lookup)
	case "/v1/flag":
		flagRedirect(w, r, lookup)
	case "/v1/map":
		mapRedirect(w, r, lookup)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func writeLookup(w http.ResponseWriter, lookup geoip.LookupResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lookup)
}

func mapRedirect(w http.ResponseWriter, r *http.Request, lookup geoip.LookupResponse) {
	http.Redirect(
		w, r, fmt.Sprintf(conf.GetString("map"), lookup.Latitude, lookup.Longitude),
		http.StatusPermanentRedirect,
	)
}
func flagRedirect(w http.ResponseWriter, r *http.Request, lookup geoip.LookupResponse) {
	http.Redirect(
		w, r, fmt.Sprintf(conf.GetString("flag"), strings.ToLower(lookup.Country)),
		http.StatusPermanentRedirect,
	)
}

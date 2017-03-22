package handler

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	conf "github.com/spf13/viper"
	"glok/geoip"
	"net"
	"net/http"
	"strings"
)

func HandlerV1(w http.ResponseWriter, r *http.Request) {
	ipString := r.FormValue("ip")
	if ipString == "" {
		ipString, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	lookup := geoip.Lookup(net.ParseIP(ipString))

	if conf.GetBool("debug") {
		log.Debug(fmt.Sprintf("Host: %s, Url: %s", r.RemoteAddr, r.URL))
	}

	switch r.URL.Path {
	case "/v1/lookup":
		writeLookup(w, lookup)
		return
	case "/v1/flag":
		flagRedirect(w, r, lookup)
		return
	case "/v1/map":
		mapRedirect(w, r, lookup)
		return
	}
	http.NotFound(w, r)
}

func writeLookup(w http.ResponseWriter, lookup geoip.LookupResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lookup)
}

func mapRedirect(w http.ResponseWriter, r *http.Request, lookup geoip.LookupResponse) {
	http.Redirect(w, r, fmt.Sprintf(conf.GetString("map"), lookup.Latitude, lookup.Longitude), 302)
}
func flagRedirect(w http.ResponseWriter, r *http.Request, lookup geoip.LookupResponse) {
	http.Redirect(w, r, fmt.Sprintf(conf.GetString("flag"), strings.ToLower(lookup.Country)), 302)
}

package handler

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"net"
)

func Base(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.Debugf("[Middlware] %s Incoming %s method from: %s ", r.URL, r.Method, ip)

	next(w, r)
}

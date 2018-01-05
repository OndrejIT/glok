package handler

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	conf "github.com/spf13/viper"
	"net"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"bytes"
	"fmt"
	"strings"
	"encoding/base64"
)

func Base(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.Debugf("[Middleware] %s Incoming %s method from: %s ", r.URL, r.Method, ip)
	if err := authCheck(r, w); err == nil {
		next(w, r)
	}
}

func authCheck(r *http.Request, w http.ResponseWriter) error{
	allowed_ips := conf.GetStringSlice("allowed_ips")
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	remote_ip := net.ParseIP(host)

	if ipInNetwork(remote_ip, allowed_ips) {
		return nil
	}

	if r.FormValue("token") != "" {
		return check_token(r, w)
	} else {
		user, password := getBasicAuth(r)
		if user != os.Getenv("API_USER") || password != os.Getenv("API_PASSWORD") {
			w.WriteHeader(403)
			log.Error("Forbidden: wrong user or password")
			return errors.New("Forbidden: wrong user or password")
		}
	}
	return nil
}

func getBasicAuth(r *http.Request) (string, string) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
            return "", ""
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		log.Error("Missing user or password")
		return "", ""
	}

	return pair[0], pair[1]
}

func check_token(r *http.Request, w http.ResponseWriter) error {
	tokenString := r.FormValue("token")
	if tokenString == "" {
		log.Debug("Missing security token.")
		w.WriteHeader(403)
		return errors.New("Missing security token")
	}

	secret_key := []byte(conf.GetString("secret_key"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret_key, nil
	})

	if token != nil && token.Valid {
		return nil
	}

	log.Error(err)
	w.WriteHeader(403)
	return errors.New("Invalid token")
}

func ipInNetwork(ip net.IP, slice []string) bool {
	for _, allowed := range slice {
		parsed := net.ParseIP(allowed)
		if parsed != nil {
			if bytes.Compare(ip, parsed) == 0 {
				return true
			}
		} else {
			_, network, err := net.ParseCIDR(allowed)
			if err != nil {
				log.Errorf("%s is not valid ip", allowed)
				return false
			}

			if network.Contains(ip) {
				return true
			}
		}
	}
	return false
}
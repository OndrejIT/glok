package geoip

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oschwald/geoip2-golang"
	conf "github.com/spf13/viper"
)

var db *geoip2.Reader

func Setup() {
	var err error
	db, err = geoip2.Open(conf.GetString("database"))
	if err != nil {
		log.Fatalf("[Geoip] %s", err)
	}
}

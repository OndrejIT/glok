package geoip

import (
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
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

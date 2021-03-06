package main

import (
	"github.com/ondrejit/glok/config"
	"github.com/ondrejit/glok/geoip"
	"github.com/ondrejit/glok/handler"
)

func init() {
	config.Setup()
	geoip.DBupdate()
}

func main() {
	handler.Start()
}

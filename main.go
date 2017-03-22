package main

import (
	"glok/config"
	"glok/geoip"
	"glok/handler"
)

func init() {
	config.Setup()
	geoip.Setup()
}

func main() {
	handler.Start()
}

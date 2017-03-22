package geoip

import (
	"bytes"
	conf "github.com/spf13/viper"
	"net"
	"testing"
)

var yamlConf = []byte(`
database: "../GeoIP2-City-Test.mmdb"
map: "http://maps.google.com?q=%f,%f"
flag: "http://www.theodora.com/flags/%s-t.gif"
`)

func init() {
	conf.SetConfigType("yaml")
	conf.ReadConfig(bytes.NewBuffer(yamlConf))
	Setup()
}

func TestLookup(t *testing.T) {
	lookup := Lookup(net.ParseIP("81.2.69.160"))
	if lookup.Country != "GB" {
		t.Error("Country")
	}
	if lookup.City != "London" {
		t.Error("City")
	}
	if lookup.Longitude != -0.0931 {
		t.Error("Longitude")
	}
	if lookup.Latitude != 51.5142 {
		t.Error("Latitude")
	}
}

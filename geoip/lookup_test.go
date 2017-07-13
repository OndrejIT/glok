package geoip

import (
	"bytes"
	conf "github.com/spf13/viper"
	"net"
	"testing"
	"net/http"
)

var yamlConf = []byte(`
map: "http://maps.google.com?q=%f,%f"
flag: "http://www.theodora.com/flags/%s-t.gif"
`)

func init() {
	conf.SetConfigType("yaml")
	conf.ReadConfig(bytes.NewBuffer(yamlConf))
	conf.Set("database", "../GeoIP2-City-Test.mmdb")
	Setup()
}

func TestLookup(t *testing.T) {

	lookup, status, _ := Lookup(net.ParseIP("81.2.69.160"))
	if status != http.StatusOK {
		t.Errorf("Bad status code: %d", status)
	}
	if lookup.Country != "GB" {
		t.Errorf("Bad country: %s", lookup.Country)
	}
	if lookup.City != "London" {
		t.Errorf("Bad city: %s", lookup.City)
	}
	if lookup.Longitude != -0.0931 {
		t.Errorf("Bad longitude %f", lookup.Longitude)
	}
	if lookup.Latitude != 51.5142 {
		t.Errorf("Bad latitude %f", lookup.Latitude)
	}
}

package api

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"bytes"
	"github.com/ondrejit/glok/geoip"
	conf "github.com/spf13/viper"
	"encoding/json"
	"io/ioutil"
)

var yamlConf = []byte(`
database: "../GeoIP2-City-Test.mmdb"
map: "http://maps.google.com?q=%f,%f"
flag: "http://www.theodora.com/flags/%s-t.gif"
`)

var ipaddr = string("81.2.69.166")

func init() {
	conf.SetConfigType("yaml")
	conf.ReadConfig(bytes.NewBuffer(yamlConf))
	geoip.Setup()
}

func TestV1BadRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/BadRequest", nil)
	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(V1)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Bad status code: %d", resp.Code)
	}

}

func TestV1Response(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/lookup", nil)
	if err != nil {
		t.Error(err)
	}

	q := req.URL.Query()
	q.Add("ip", ipaddr)
	req.URL.RawQuery = q.Encode()

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(V1)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Bad status code: %d", resp.Code)
	}

	lookup := geoip.LookupResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	json.Unmarshal(body, &lookup)

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

func TestV1Flag(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/flag", nil)
	if err != nil {
		t.Error(err)
	}

	q := req.URL.Query()
	q.Add("ip", ipaddr)
	req.URL.RawQuery = q.Encode()

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(V1)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusPermanentRedirect {
		t.Errorf("Bad status code: %d", resp.Code)
	}

	flagRedirect := "<a href=\"http://www.theodora.com/flags/gb-t.gif\">Permanent Redirect</a>.\n\n"
	if resp.Body.String() != flagRedirect {
		t.Errorf("Bad redirect: %s", resp.Body.String())
	}
}

func TestV1Map(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/map", nil)
	if err != nil {
		t.Error(err)
	}

	q := req.URL.Query()
	q.Add("ip", ipaddr)
	req.URL.RawQuery = q.Encode()

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(V1)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusPermanentRedirect {
		t.Errorf("Bad status code: %d", resp.Code)
	}

	mapRedirect := "<a href=\"http://maps.google.com?q=51.514200,-0.093100\">Permanent Redirect</a>.\n\n"
	if resp.Body.String() != mapRedirect {
		t.Errorf("Bad redirect: %s", resp.Body.String())
	}
}

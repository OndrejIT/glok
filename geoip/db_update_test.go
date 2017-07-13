package geoip

import (
	"testing"
	conf "github.com/spf13/viper"
	"os"
)


func TestUpdateMissingDb(t *testing.T) {
	conf.Set("database", "testDB.mmdb")
	conf.Set("license", "000000000000")
	conf.Set("uid", "0")
	conf.Set("product_id", "GeoLite2-Country")


	err := update_db()
	if err != nil {
		t.Error(err)
	}

	os.Remove(conf.GetString("database"))
}

func TestUpdateExistDb(t *testing.T) {
	// vytvorim prazdny soubor s nejakou md5
	os.Create(conf.GetString("database"))

	err := update_db()
	if err != nil {
		t.Error(err)
	}
	// stahnuta aktualni db
}

func TestUpdateNewestDb(t *testing.T) {
	// v predchozim testu jsem stahnul aktualni db => update_db na to prijde a nebude stahovat novou
	err := update_db()
	if err != nil {
		t.Error(err)
	}

	// a nakonec uklidim
	os.Remove(conf.GetString("database"))
}

package config

import (
	flag "github.com/spf13/pflag"
	conf "github.com/spf13/viper"
)

func FlagParser() {
	flag.IntP("port", "p", 8888, "Set server port.")
	flag.String("ip", "0.0.0.0", "Set server ip.")
	flag.String("db", "GeoIP2-City.mmdb", "Set database path.")
	flag.StringP("flag", "f", "http://www.theodora.com/flags/%s-t.gif", "Set flag url.")
	flag.StringP("map", "m", "http://maps.google.com?q=%f,%f", "Set map url.")
	flag.StringP("config", "c", ".", "Set config path.")
	flag.BoolP("debug", "d", false, "Enable debug mode.")
	flag.IntP("interval", "i", 24, "DB update interval")
	flag.Parse()
	flagToConfig()
}

func flagToConfig() {
	conf.BindPFlag("ip", flag.Lookup("ip"))
	conf.BindPFlag("port", flag.Lookup("port"))
	conf.BindPFlag("database", flag.Lookup("db"))
	conf.BindPFlag("flag", flag.Lookup("flag"))
	conf.BindPFlag("map", flag.Lookup("map"))
	conf.BindPFlag("config", flag.Lookup("config"))
	conf.BindPFlag("debug", flag.Lookup("debug"))
	conf.BindPFlag("interval", flag.Lookup("interval"))
}

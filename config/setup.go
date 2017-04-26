package config

import (
	log "github.com/Sirupsen/logrus"
	conf "github.com/spf13/viper"
)

func Setup() {
	FlagParser()

	if conf.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
		text_formatter := new(log.TextFormatter)
		text_formatter.TimestampFormat = "2006-01-02 15:04:05"
		log.SetFormatter(text_formatter)
		text_formatter.FullTimestamp = true
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	conf.SetConfigName("glok")
	conf.AddConfigPath(conf.GetString("config"))

	conf.SetEnvPrefix("glok")
	conf.AutomaticEnv()

	if err := conf.ReadInConfig(); err != nil {
		log.Fatal("[Config] Fatal error config file: ", err)
	}

}

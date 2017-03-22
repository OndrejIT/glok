package handler

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	conf "github.com/spf13/viper"
	"net/http"
)

func Start() {
	hostPort := fmt.Sprintf("0.0.0.0:%d", conf.GetInt("port"))
	log.Info("[Main] Initiating server listening at ", hostPort)
	http.HandleFunc("/v1/", HandlerV1)
	log.Fatal(http.ListenAndServe(hostPort, nil))
}

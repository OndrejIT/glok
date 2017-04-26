package handler

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	conf "github.com/spf13/viper"
	"net/http"
	"github.com/urfave/negroni"
	"github.com/ondrejit/glok/api"
)

func Start() {
	hostPort := fmt.Sprintf("0.0.0.0:%d", conf.GetInt("port"))
	log.Info("[Main] Initiating server listening at ", hostPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/", api.V1)

	n := negroni.New(negroni.HandlerFunc(Base))
	n.UseHandler(mux)

	log.Fatal(http.ListenAndServe(hostPort, n))
}

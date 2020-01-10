package handler

import (
	"fmt"
	"net/http"

	"github.com/ondrejit/glok/api"
	log "github.com/sirupsen/logrus"
	conf "github.com/spf13/viper"
	"github.com/urfave/negroni"
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

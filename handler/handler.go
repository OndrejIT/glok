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
	address := fmt.Sprintf("%s:%d", conf.GetString("ip"), conf.GetInt("port"))
	log.Info("[Main] Initiating server listening at ", address)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/", api.V1)

	n := negroni.New(negroni.HandlerFunc(Base))
	n.UseHandler(mux)

	log.Fatal(http.ListenAndServe(address, n))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	appconfig "github.com/lelledaniele/upaygo/config"
	apprestintentcancel "github.com/lelledaniele/upaygo/controller/rest/intent/cancel"
	apprestintentcapture "github.com/lelledaniele/upaygo/controller/rest/intent/capture"
	apprestintentconfirm "github.com/lelledaniele/upaygo/controller/rest/intent/confirm"
	apprestintentcreate "github.com/lelledaniele/upaygo/controller/rest/intent/create"
	apprestintentget "github.com/lelledaniele/upaygo/controller/rest/intent/get"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path for config file")
	flag.Parse()

	if configFile == "" {
		log.Fatal("Flag 'config' for configuration file path is required")
	}
}

// @title uPayment in GO
// @version 1.0.0
// @description Microservice to manage payment
// @license.name MIT
func main() {
	fc, e := os.Open(configFile)
	if e != nil {
		log.Fatal(fmt.Sprintf("Impossible to open configuration file: %v\n", e))
	}
	defer fc.Close()

	e = appconfig.ImportConfig(fc)
	if e != nil {
		log.Fatal(fmt.Sprintf("Error durring file config import: %v\n", e))
	}

	s := appconfig.GetServerConfig()

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc(apprestintentget.URL, apprestintentget.Handler).Methods(apprestintentget.Method)
	r.HandleFunc(apprestintentcreate.URL, apprestintentcreate.Handler).Methods(apprestintentcreate.Method)
	r.HandleFunc(apprestintentconfirm.URL, apprestintentconfirm.Handler).Methods(apprestintentconfirm.Method)
	r.HandleFunc(apprestintentcapture.URL, apprestintentcapture.Handler).Methods(apprestintentcapture.Method)
	r.HandleFunc(apprestintentcancel.URL, apprestintentcancel.Handler).Methods(apprestintentcancel.Method)

	log.Fatal(http.ListenAndServe(":"+s.GetPort(), r))
}

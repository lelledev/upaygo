package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	appconfig "github.com/lelledev/upaygo/config"
	apprestintentcancel "github.com/lelledev/upaygo/controller/rest/intent/cancel"
	apprestintentcapture "github.com/lelledev/upaygo/controller/rest/intent/capture"
	apprestintentconfirm "github.com/lelledev/upaygo/controller/rest/intent/confirm"
	apprestintentcreate "github.com/lelledev/upaygo/controller/rest/intent/create"
	apprestintentget "github.com/lelledev/upaygo/controller/rest/intent/get"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"

	_ "github.com/lelledev/upaygo/docs"
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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fc, e := os.Open(configFile) //nolint:gosec // configFile is an operator-supplied CLI flag, not user input
	if e != nil {
		return fmt.Errorf("impossible to open configuration file: %w", e)
	}
	defer func() {
		_ = fc.Close()
	}()

	if e = appconfig.ImportConfig(fc); e != nil {
		return fmt.Errorf("error during file config import: %w", e)
	}

	s := appconfig.GetServerConfig()

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc(apprestintentget.URL, apprestintentget.Handler).Methods(apprestintentget.Method)
	r.HandleFunc(apprestintentcreate.URL, apprestintentcreate.Handler).Methods(apprestintentcreate.Method)
	r.HandleFunc(apprestintentconfirm.URL, apprestintentconfirm.Handler).Methods(apprestintentconfirm.Method)
	r.HandleFunc(apprestintentcapture.URL, apprestintentcapture.Handler).Methods(apprestintentcapture.Method)
	r.HandleFunc(apprestintentcancel.URL, apprestintentcancel.Handler).Methods(apprestintentcancel.Method)

	srv := &http.Server{
		Addr:              ":" + s.GetPort(),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return srv.ListenAndServe()
}

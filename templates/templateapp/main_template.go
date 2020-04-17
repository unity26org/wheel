package templates

var MainPath = []string{"main.go"}

var MainContent = `package main

import (
	"flag"
	"net/http"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/log"
	"{{ .AppRepository }}/config"
	"{{ .AppRepository }}/db/schema"
	"{{ .AppRepository }}/routes"
)

func main() {
	var mode string
	var port string
	var host string

	flag.StringVar(&mode, "mode", "server", "run mode (options: server/migrate)")
	flag.StringVar(&host, "host", "localhost", "http server host")
	flag.StringVar(&port, "port", "8081", "http server port")
	flag.Parse()

	log.Info.Println("starting app", config.App.AppName)

	model.Connect()

	if mode == "migrate" {
    schema.Migrate()
  } else if mode == "rollback" {
    schema.Rollback()	
  } else if mode == "s" || mode == "server" {
		log.Fatal.Println(http.ListenAndServe(host+":"+port, routes.Routes(host, port)))
	} else {
		log.Fatal.Println("invalid run mode, please, use \"--help\" for more details")
	}
}`

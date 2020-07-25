package templates

var MainPath = []string{"main.go"}

var MainContent = `package main

import (
  "github.com/adilsonchacon/sargo"
	"net/http"
  "os"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/log"
	"{{ .AppRepository }}/config"
	"{{ .AppRepository }}/db/schema"
	"{{ .AppRepository }}/routes"
)

func main() {
	sargo.Set("mode", "m", "server", "run mode (options: server/migrate). Default value is \"server\"")
  sargo.Set("binding", "b", "localhost", "http server IP. Default value is \"localhost\"")
  sargo.Set("port", "p", 8081, "http server port. Default value is \"8081\"")

  if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
  	sargo.PrintHelpAndExit()
  }
  
  mode, _ := sargo.Get("mode")
  binding, _ := sargo.Get("binding")
  port, _ := sargo.Get("port")

	log.Info.Println("starting app", config.App.AppName)

	if mode == "migrate" {
  	model.Connect()
    schema.Migrate()
  } else if mode == "rollback" {
  	model.Connect()
    schema.Rollback()	
  } else if mode == "s" || mode == "server" {
  	model.Connect()
		log.Fatal.Println(http.ListenAndServe(binding+":"+port, routes.Routes(binding, port)))
	} else {
		log.Fatal.Println("invalid run mode, please, use \"--help\" for more details")
	}
}`

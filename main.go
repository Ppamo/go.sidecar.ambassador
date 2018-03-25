package main

import (
	"github.com/Ppamo/go.sidecar.ambassador/config"
	"github.com/Ppamo/go.sidecar.ambassador/server"
	"log"
)

func main() {
	e := config.LoadConfig("config.json")
	if e != nil {
		log.Fatalf("- Error loading conf\n%v\n", e)
		panic(e)
	}
	e = server.StartServer()
	if e != nil {
		log.Fatalf("- Error iniciando servicio de traduccion\n%v\n", e)
	}
}

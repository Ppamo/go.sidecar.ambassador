package main

import (
	"github.com/Ppamo/go.sidecar.ambassador/config"
	"github.com/Ppamo/go.sidecar.ambassador/rules"
	"github.com/Ppamo/go.sidecar.ambassador/server"
	"log"
)

func main() {
	e := config.LoadConfig("config.json")
	if e != nil {
		log.Panicf("- ERROR: Fail to load configuration\n%v\n", e)
	}
	e = rules.LoadRules()
	if e != nil {
		log.Panicf("- ERROR: Fail to load rules from host\n%v\n", e)
	}
	e = server.StartServer()
	if e != nil {
		log.Panicf("- ERROR: Starting validation service\n%v\n", e)
	}
}

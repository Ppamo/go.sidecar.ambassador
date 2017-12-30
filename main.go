package main

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/config"
	"github.com/Ppamo/go.sidecar.ambassador/server"
)

func main() {
	config, e := config.LoadConfig("config/config_test.json")
	if e != nil {
		fmt.Printf("Error generating conf\n%v\n", e)
	}
	server.StartServer(config)
	fmt.Printf("server port: %v\n", config.Server.Port)
}

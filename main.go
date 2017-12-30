package main

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/config"
)

func main() {
	config, e := config.LoadConfig("config/main.json")
	if e != nil {
		fmt.Printf("Error generating conf\n%v\n", e)
	}
	fmt.Printf("server port: %i\n", config.Server.Port)
}

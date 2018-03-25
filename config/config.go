package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Server ServerConfig `json:"server"`
	Host   HostConfig   `json:"parent"`
}

type ServerConfig struct {
	Host string `json: "host"`
	Port int    `json:"port"`
}

type HostConfig struct {
	Destination    string `json:"destination"`
	UrlPrefix      string `json:"urlPrefix"`
	ServiceInfoUrl string `json:"serviceInfoUrl"`
}

func LoadConfig(configPath string) (Configuration, error) {
	var config Configuration
	file, e := os.Open(configPath)
	defer file.Close()
	if e != nil {
		log.Fatalf("- Error reading conf\n%v\n", e)
		panic(e)
	}
	parser := json.NewDecoder(file)
	e = parser.Decode(&config)
	return config, e
}

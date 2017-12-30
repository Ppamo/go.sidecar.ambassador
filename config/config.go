package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	*Server `json:"server"`
	*Parent `json:"parent"`
}

type Server struct {
	Host string `json: "host"`
	Port int    `json:"port"`
}

type Parent struct {
	Destination    string `json:"destination"`
	UrlPrefix      string `json:"urlPrefix"`
	ServiceInfoUrl string `json:"serviceInfoUrl"`
}

func LoadConfig(configPath string) (Configuration, error) {
	var config Configuration
	file, e := os.Open(configPath)
	defer file.Close()
	if e != nil {
		return config, e
	}
	parser := json.NewDecoder(file)
	e = parser.Decode(&config)
	return config, e
}

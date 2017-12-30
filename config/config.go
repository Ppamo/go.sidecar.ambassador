package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	*Server `json:"server"`
	*Parent `json:"parent"`
}

type Server struct {
	Port int `json:"port"`
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
		fmt.Printf("Could not open file %s\n%v\n", configPath, e)
		return config, e
	}
	parser := json.NewDecoder(file)
	e = parser.Decode(&config)
	return config, e
}

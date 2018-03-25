package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

type Configuration struct {
	Server ServerConfig `json:"server"`
	Host   HostConfig   `json:"host"`
}

type ServerConfig struct {
	Host string `json:"host" env:"SERVERHOST"`
	Port int    `json:"port" env:"SERVERPORT"`
}

type HostConfig struct {
	Destination    string `json:"destination" env:"DESTINATION"`
	UrlPrefix      string `json:"urlPrefix" env:"URLPREFIX"`
	ServiceInfoUrl string `json:"serviceInfoUrl" env:"SERVICEINFO"`
}

func setValues(item interface{}) {
	var stringVal string
	rt := reflect.TypeOf(item)
	rv := reflect.ValueOf(item)

	for i := 0; i < rt.NumField(); i++ {
		key := rt.Field(i).Tag.Get("env")
		value := rv.Field(i).Interface()
		if rv.Field(i).Kind() == reflect.Struct {
			setValues(value)
			continue
		}
		if len(key) > 0 {
			switch rv.Field(i).Kind() {
			case reflect.String:
				stringVal = fmt.Sprintf("%s", value.(string))
			case reflect.Int32, reflect.Int, reflect.Int64:
				stringVal = fmt.Sprintf("%d", value.(int))
			case reflect.Float32:
				stringVal = fmt.Sprintf("%f", value.(float32))
			case reflect.Float64:
				stringVal = fmt.Sprintf("%f", value.(float64))
			case reflect.Bool:
				stringVal = fmt.Sprintf("%b", value.(bool))
			}
			_, ok := os.LookupEnv(key)
			if !ok {
				log.Printf("+ Seting %s=%s\n", key, stringVal)
				os.Setenv(key, value.(string))
			} else {
				log.Printf("+ Skipping %s=%s\n", key, stringVal)
			}
		}
	}
}

func setConfigValues(config *Configuration) error {
	setValues(config.Server)
	setValues(config.Host)
	return nil
}

func LoadConfig(configPath string) error {
	var config Configuration
	file, e := os.Open(configPath)
	defer file.Close()
	if e != nil {
		log.Fatalf("- Error reading conf\n%v\n", e)
		panic(e)
	}
	parser := json.NewDecoder(file)
	e = parser.Decode(&config)
	if e != nil {
		log.Fatalf("- Error parsing config json\n%v\n", e)
		panic(e)
	}
	e = setConfigValues(&config)
	return e
}

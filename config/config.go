package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

type HostProperties struct {
	Items []HostProperty `json:"properties"`
}

type HostProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Configuration struct {
	Server ServerConfig `json:"server"`
	Host   HostConfig   `json:"host"`
}

type ServerConfig struct {
	Host         string `json:"host" env:"SERVERHOST"`
	Port         int    `json:"port" env:"SERVERPORT"`
	RequestRetry int    `json:"requestRetry" env:"REQUESTRETRY"`
}

type HostConfig struct {
	Destination       string `json:"destination" env:"DESTINATION"`
	UrlPrefix         string `json:"urlPrefix" env:"URLPREFIX"`
	HostRulesURL      string `json:"hostRulesUrl" env:"HOSTRULESURL"`
	HostPropertiesURL string `json:"hostPropertiesUrl" env:"HOSTPROPERTIESURL"`
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
			envValue, ok := os.LookupEnv(key)
			if !ok {
				log.Printf("+ Seting %s=%s\n", key, stringVal)
				os.Setenv(key, stringVal)
			} else {
				log.Printf("+ Skipping %s:%s\n", key, envValue)
			}
		}
	}
}

func loadHostProperties() error {
	var err error
	var response *http.Response
	url := fmt.Sprintf("%s%s", os.Getenv("DESTINATION"), os.Getenv("HOSTPROPERTIESURL"))
	retry, _ := strconv.Atoi(os.Getenv("REQUESTRETRY"))
	for i := 0; i < retry; i++ {
		log.Printf("+ Getting properties, attempt #%d\n", i+1)
		response, err = http.Get(url)
		if err != nil {
			response = nil
			log.Printf("- ERROR: Fail to get properties from:\n%s\n%v", url, err)
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)
		var properties HostProperties
		err = decoder.Decode(&properties)
		if err != nil {
			log.Printf("- ERROR: Could not decode properties\n%v\n", err)
			return err
		}
		for _, item := range properties.Items {
			log.Printf("+ Setting %s=%s\n", item.Key, item.Value)
			os.Setenv(item.Key, item.Value)
		}
	} else {
		log.Printf("- ERROR: Could not load properties from %s\n", url)
		return errors.New("Failed to load properties")
	}

	return nil
}

func LoadConfig(configPath string) error {
	var config Configuration
	file, e := os.Open(configPath)
	defer file.Close()
	if e != nil {
		log.Panicf("- Error reading conf\n%v\n", e)
	}
	parser := json.NewDecoder(file)
	e = parser.Decode(&config)
	if e != nil {
		log.Panicf("- Error parsing config json\n%v\n", e)
	}
	setValues(config)
	e = loadHostProperties()
	return e
}

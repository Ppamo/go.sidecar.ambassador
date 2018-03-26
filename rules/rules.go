package rules

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Rules struct {
	Operations []Operation `json:"enabled"`
}

type Operation struct {
	Description string                 `json:"description"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Params      map[string]interface{} `json:"params"`
	Body        map[string]interface{} `json:"body"`
}

var rules map[string]Operation

func loadHostRules() (*Rules, error) {
	var err error
	var operations *Rules
	var response *http.Response
	url := fmt.Sprintf("%s%s", os.Getenv("DESTINATION"), os.Getenv("HOSTRULESURL"))
	retry, _ := strconv.Atoi(os.Getenv("REQUESTRETRY"))
	for i := 0; i < retry; i++ {
		log.Printf("+ Getting rules, attempt #%d\n", i+1)
		response, err = http.Get(url)
		if err != nil {
			response = nil
			log.Printf("- ERROR: Fail to get response from:\n%s\n%v", url, err)
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)
		operations = new(Rules)
		err = decoder.Decode(&operations)
		if err != nil {
			log.Printf("- ERROR: Could not decode response\n%v\n", err)
			return nil, err
		}
	}
	return operations, nil
}

func getMapKey(method string, url string) string {
	return fmt.Sprintf("%s::%s", strings.ToLower(method), strings.ToLower(url))
}

func mapRules(hostRules *Rules) (map[string]Operation, error) {
	var mapkey string
	rules = make(map[string]Operation)
	for _, item := range hostRules.Operations {
		mapkey = getMapKey(item.Method, item.Path)
		rules[mapkey] = item
	}
	return rules, nil
}

func GetRule(method string, url string) *Operation {
	log.Printf("+ Geting rule %s::%s", method, url)
	operation, ok := rules[getMapKey(method, url)]
	var copy *Operation
	if ok {
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		decoder := gob.NewDecoder(&buffer)
		err := encoder.Encode(operation)
		if err != nil {
			log.Printf("- ERROR: Could not encode rule\n%v", err)
			return nil
		}
		copy = new(Operation)
		err = decoder.Decode(&copy)
		if err != nil {
			log.Printf("- ERROR: Could not dencode rule\n%v", err)
			return nil
		}
		return copy

	}
	log.Printf("- ERROR: Rule not found\n")
	return nil
}

func registerGob() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}

func LoadRules() error {
	registerGob()
	hostRules, err := loadHostRules()
	if err != nil {
		log.Panicf("- ERROR: Failed to load rules\n%v\n", err)
	}
	rules, err = mapRules(hostRules)
	return err
}

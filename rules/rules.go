package rules

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/structs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var rules map[string]structs.Operation

func loadHostRules() (*structs.Rules, error) {
	var err error
	var rules *structs.Rules
	var response *http.Response
	url := fmt.Sprintf("%s%s", os.Getenv("DESTINATION"), os.Getenv("HOSTRULESURL"))
	retry, _ := strconv.Atoi(os.Getenv("REQUESTRETRY"))
	for i := 0; i < retry; i++ {
		log.Printf("+ Loading rules, attempt #%d\n", i+1)
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
		rules = new(structs.Rules)
		err = decoder.Decode(&rules)
		if err != nil {
			log.Printf("- ERROR: Could not decode response\n%v\n", err)
			return nil, err
		}
	}
	return rules, nil
}

func getMapKey(method string, url string) string {
	return fmt.Sprintf("%s::%s", strings.ToLower(method), strings.ToLower(url))
}

func mapRules(hostRules *structs.Rules) (map[string]structs.Operation, error) {
	var mapkey string
	var err error
	rules = make(map[string]structs.Operation)
	for _, item := range hostRules.Operations {
		mapkey = getMapKey(item.Method, item.Path)
		// item.ParamsSchema, err = validator.GetCompiledSchema(item.Params)
		// item.BodySchema, err = validator.GetCompiledSchema(item.Body)
		if err != nil {
			return nil, err
		}
		rules[mapkey] = item
	}
	return rules, nil
}

func GetRule(method string, url string) *structs.Operation {
	log.Printf("+ Getting rule %s::%s", method, url)
	operation, ok := rules[getMapKey(method, url)]
	var copy *structs.Operation
	if ok {
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		decoder := gob.NewDecoder(&buffer)
		err := encoder.Encode(operation)
		if err != nil {
			log.Printf("- ERROR: Could not encode rule\n%v", err)
			return nil
		}
		copy = new(structs.Operation)
		err = decoder.Decode(&copy)
		if err != nil {
			log.Printf("- ERROR: Could not dencode rule\n%v", err)
			return nil
		}
		return copy

	}
	log.Printf("- ERROR: Operation not found\n")
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

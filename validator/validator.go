package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Ppamo/go.sidecar.ambassador/structs"
	"github.com/xeipuuv/gojsonschema"
	"log"
	"net/http"
)

func GetCompiledSchema(schema map[string]interface{}) (map[string]interface{}, error) {
	if schema == nil {
		return nil, nil
	}
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(schema)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func ValidateParams(rule *structs.Operation, request *http.Request) error {
	var jsonValues, jsonSchema []byte
	var err error
	if rule.Params == nil || len(rule.Params) == 0 {
		jsonSchema = []byte("{\"additionalProperties\": false}")
	} else {
		jsonSchema, err = json.Marshal(rule.Params)
		if err != nil {
			return err
		}
	}
	jsonValues, err = json.Marshal(request.URL.Query())
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewStringLoader(string(jsonSchema))
	paramsLoader := gojsonschema.NewStringLoader(string(jsonValues))
	result, err := gojsonschema.Validate(schemaLoader, paramsLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		for _, description := range result.Errors() {
			log.Printf("- VALIDATION ERROR: %s\n", description)
		}
		return errors.New("Schema validation errors")
	}
	return nil
}

func ValidateBody(rule *structs.Operation, request *http.Request) error {
	return nil
}

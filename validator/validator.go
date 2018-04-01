package validator

import (
	"encoding/json"
	"errors"
	"github.com/Ppamo/go.sidecar.ambassador/structs"
	gjs "github.com/xeipuuv/gojsonschema"
	"log"
	"net/http"
)

func GetJSONSchemas(rule structs.Operation) (string, string, error) {
	var paramsSchema, bodySchema string

	jsonSchema, err := json.Marshal(rule.Params)
	if err != nil || len(jsonSchema) == 4 {
		jsonSchema = []byte("{\"additionalProperties\": false}")
	}
	paramsSchema = string(jsonSchema)

	jsonSchema, err = json.Marshal(rule.Body)
	if err != nil || len(jsonSchema) == 4 {
		jsonSchema = []byte("{\"additionalProperties\": false}")
	}
	bodySchema = string(jsonSchema)

	return paramsSchema, bodySchema, nil
}

func validate(schema string, code string) error {
	schemaLoader := gjs.NewStringLoader(schema)
	paramsLoader := gjs.NewStringLoader(code)
	result, err := gjs.Validate(schemaLoader, paramsLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		for _, description := range result.Errors() {
			log.Printf("- VALIDATION ERROR: %s\n", description)
		}
		return errors.New("Schema validation failed!")
	}
	return nil
}

func ValidateParams(rule *structs.Operation, request *http.Request) error {
	valuesCode, err := json.Marshal(request.URL.Query())
	if err != nil {
		return err
	}
	return validate(rule.ParamsCode, string(valuesCode))
}

func ValidateBody(rule *structs.Operation, request *http.Request) error {
	valuesCode, err := json.Marshal(request.Body)
	if err != nil {
		return err
	}
	return validate(rule.BodyCode, string(valuesCode))
}

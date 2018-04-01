package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/rules"
	"github.com/Ppamo/go.sidecar.ambassador/structs"
	"github.com/Ppamo/go.sidecar.ambassador/utils"
	"github.com/Ppamo/go.sidecar.ambassador/validator"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	filepath = "rules/petstore.swagger.api.json"
)

var server *http.Server
var apiRules structs.Rules

func getErrorResponse(code int, message string) string {
	response := fmt.Sprintf(`{
	"httpCode": %d,
	"httpMessage": "%s"
	"moreInformation": "%[1]d: %[2]s"
}`, code, message)
	return response
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	urlPrefix := utils.Getenv("URLPREFIX", "")
	if !strings.HasPrefix(r.URL.Path, urlPrefix) {
		http.Error(w, getErrorResponse(400, "Bad Request"), http.StatusBadRequest)
		log.Printf("- Unauthorized!\nInvalid request prefix: %s != %s", r.URL.Path, urlPrefix)
		return
	}
	url := strings.TrimPrefix(r.URL.Path, urlPrefix)
	rule := rules.GetRule(r.Method, url)
	if rule == nil {
		http.Error(w, getErrorResponse(403, "Forbidden"), http.StatusForbidden)
		log.Printf("- Unauthorized!\nNo rule found for request %s::%s", r.Method, url)
		return
	}
	log.Printf("+ Operation: %s", rule.Description)
	err := validator.ValidateParams(rule, r)
	if err != nil {
		http.Error(w, getErrorResponse(400, "Bad Request"), http.StatusBadRequest)
		log.Printf("- Unauthorized!\n%v", err)
		return
	}
	err = validator.ValidateBody(rule, r)
	if err != nil {
		http.Error(w, getErrorResponse(400, "Bad Request"), http.StatusBadRequest)
		log.Printf("- Unauthorized!\n%v", err)
		return
	}
	log.Printf("+ Autorized!")
	url = fmt.Sprintf("%s/%s", utils.Getenv("DESTINATION", ""), url)
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, getErrorResponse(404, "Not Found"), http.StatusNotFound)
		log.Printf("- ERROR: Failed request to url\n%v", url, err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, getErrorResponse(500, "Internal server error"), http.StatusInternalServerError)
		log.Printf("- ERROR: Failed to read respoonse body\n%v", err)
		return
	}
	fmt.Fprintf(w, string(body))
}

func StartServer() error {
	address := fmt.Sprintf("%s:%s",
		utils.Getenv("SERVERHOST", "0.0.0.0"),
		utils.Getenv("SERVERPORT", "8080"))
	log.Printf("* Starting server at %s\n", address)
	server = &http.Server{Addr: address}
	http.HandleFunc("/", requestHandler)
	err := server.ListenAndServe()
	return err
}

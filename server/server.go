package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/rules"
	"github.com/Ppamo/go.sidecar.ambassador/utils"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	filepath = "rules/petstore.swagger.api.json"
)

var server *http.Server
var apiRules rules.Rules

func getErrorResponse(code int, message string) string {
	response := fmt.Sprintf(`{
	"httpCode": %d,
	"httpMessage": "%s"
	"moreInformation": "%[1]d: %[2]s"
}`, code, message)
	return response
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	rule := rules.GetRule(r.Method, r.URL.Path)
	if rule == nil {
		http.Error(w, getErrorResponse(403, "Forbidden"),
			http.StatusForbidden)
		log.Printf("- Unauthorized!")
		return
	}
	log.Printf("+ Operation: %s", rule.Description)
	url := fmt.Sprintf("%s/%s", utils.Getenv("DESTINATION", ""), r.URL.RequestURI())
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, getErrorResponse(404, "Not Found"),
			http.StatusNotFound)
		log.Printf("- ERROR: Failed request to url\n%v", url, err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, getErrorResponse(500, "Internal server error"),
			http.StatusInternalServerError)
		log.Printf("- ERROR: Failed to read respoonse body\n%v", err)
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

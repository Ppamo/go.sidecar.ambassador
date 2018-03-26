package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/rules"
	"github.com/Ppamo/go.sidecar.ambassador/utils"
	"log"
	"net/http"
)

const (
	filepath = "rules/petstore.swagger.api.json"
)

var server *http.Server
var apiRules rules.Rules

func getErrorResponse() string {
	response := fmt.Sprintf(`{
	"httpCode": %d,
	"httpMessage": "%s"
	"moreInformation": "%s"
}`, 500, "Internal server error", "500: Internal server error")
	return response
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	rule := rules.GetRule(r.Method, r.URL.Path)
	log.Printf("=> %s\n", rule.Description)
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
	/*
		err := rules.LoadRules(r.Method, r.URL.Path)
		if err != nil {
			http.Error(w, getErrorResponse(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "{\"status\": \"ok\"}")
		}
	*/
}

func StartServer() error {
	serverhost := utils.Getenv("SERVERHOST", "0.0.0.0")
	serverport := utils.Getenv("SERVERPORT", "8080")
	address := fmt.Sprintf("%s:%s", serverhost, serverport)
	log.Printf("* Starting server at %s\n", address)
	server = &http.Server{Addr: address}
	http.HandleFunc("/", requestHandler)
	err := server.ListenAndServe()
	return err
}

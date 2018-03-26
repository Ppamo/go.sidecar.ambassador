package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/rules"
	utils "github.com/Ppamo/go.sidecar.ambassador/utils"
	"log"
	"net/http"
)

const (
	filepath = "rules/petstore.swagger.api.json"
)

var server *http.Server
var apiRules rules.Rules

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

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if rules.IsEmpty() {
		_, e := rules.LoadRules(filepath)
		if e != nil {
			log.Fatalf("=> Error loading rules from %s\n%v\n", filepath, e)
		} else {
			log.Printf("==> rules loaded!!")
		}

	}
	fmt.Fprintf(w, "Hello my world!\n")
	fmt.Printf("=> Handling request\n")
	r.ParseForm()
	fmt.Printf("method: %s\n", r.Method)
	fmt.Printf("host: %s\n", r.Host)
	fmt.Printf("url: %s\n", r.URL.Path)
	fmt.Printf("query: %v\n", r.URL.Query())
	fmt.Printf("postform: %s\n", r.PostForm)
}

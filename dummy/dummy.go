package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var server *http.Server

func main() {
	log.Printf("+ Starting dummy process at port 8081\n")
	server = &http.Server{Addr: "0.0.0.0:8081"}
	http.HandleFunc("/", requestHandler)
	e := server.ListenAndServe()
	if e != nil {
		log.Panicf("- Error launching the server\n%v\n", e)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=validation") {
		log.Printf("+ Rules request")
		filepath := "/validation.json"
		data, e := ioutil.ReadFile(filepath)
		if e != nil {
			log.Panicf("error: loading file %s\n%v\n", filepath, e)
		}
		fmt.Fprintf(w, "%s", data)
	} else if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=properties") {
		log.Printf("+ Properties request")
		filepath := "/properties.json"
		data, e := ioutil.ReadFile(filepath)
		if e != nil {
			log.Panicf("error: loading file %s\n%v\n", filepath, e)
		}
		fmt.Fprintf(w, "%s", data)
	} else if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=quit") {
		log.Printf("+ Quit request")
		server.Shutdown(context.Background())
	} else {
		log.Printf("+ Default request: %s", r.RequestURI)
		fmt.Fprintf(w, "{\"hello\": \"world\"}\n")
	}
}

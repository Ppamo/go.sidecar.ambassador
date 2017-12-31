package main

// curl -v http://localhost:8081/serviceInfo?item=validation

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var server *http.Server

func main() {
	fmt.Printf("=> starting parent process\n")
	server = &http.Server{Addr: "localhost:8081"}
	http.HandleFunc("/", requestHandler)
	e := server.ListenAndServe()
	if e != nil {
		fmt.Printf("Error launching the server\n%v\n", e)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=validation") {
		filepath := os.Args[1]
		data, e := ioutil.ReadFile(filepath)
		if e != nil {
			fmt.Printf("error: loading file %s\n%v\n", filepath, e)
			panic(e)
		}
		fmt.Printf("=> Returning data from file %s\n", filepath)
		fmt.Fprintf(w, "%s", data)
	} else if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=quit") {
		fmt.Printf("=> quiting!")
		server.Shutdown(nil)
	} else {
		fmt.Printf("=> call with path %s\n", r.RequestURI)
		fmt.Fprintf(w, "{\"hello\": \"world\"}\n")
		fmt.Printf("=> Returning default response\n")
	}
}

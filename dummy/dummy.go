package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var server *http.Server

func main() {
	fmt.Printf("=> starting chad process at port 8081\n")
	server = &http.Server{Addr: "0.0.0.0:8081"}
	http.HandleFunc("/", requestHandler)
	e := server.ListenAndServe()
	if e != nil {
		fmt.Printf("- Error launching the server\n%v\n", e)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "get") && strings.EqualFold(r.RequestURI, "/"){
		fmt.Fprintf(w, "{\"status\": \"ok\"}")
	} else if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=validation") {
		filepath := "validation.rules.json"
		data, e := ioutil.ReadFile(filepath)
		if e != nil {
			fmt.Printf("error: loading file %s\n%v\n", filepath, e)
			panic(e)
		}
		fmt.Fprintf(w, "%s", data)
	} else if strings.EqualFold(r.Method, "get") &&
		strings.EqualFold(r.RequestURI, "/serviceInfo?item=quit") {
		fmt.Printf("+ quiting!")
		server.Shutdown(nil)
	} else {
		fmt.Printf("- call with path %s\n", r.RequestURI)
		fmt.Fprintf(w, "{\"hello\": \"world\"}\n")
		fmt.Printf("- Returning default response\n")
	}
}

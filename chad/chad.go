package main

import (
	"fmt"
	"net/http"
)

var server *http.Server

func main() {
	fmt.Printf("=> starting chad process at port 8081\n")
	server = &http.Server{Addr: "localhost:8081"}
	http.HandleFunc("/", requestHandler)
	e := server.ListenAndServe()
	if e != nil {
		fmt.Printf("Error starting the server\n%v\n", e)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("=> ok!\n")
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

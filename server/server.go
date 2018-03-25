package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/config"
	"log"
	"net/http"
	"os"
)

var server *http.Server
var serverConfig config.Configuration

func StartServer(config config.Configuration) error {
	port, ok := os.LookupEnv("serverport")
	if !ok {
		port = fmt.Sprintf("%d", config.Server.Port)
	}
	address := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("* Starting server at %s\n", address)
	server = &http.Server{Addr: address}
	http.HandleFunc("/", requestHandler)
	err := server.ListenAndServe()
	return err
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello my world!\n")
	fmt.Printf("=> Handling request\n")
	r.ParseForm()
	fmt.Printf("method: %s\n", r.Method)
	fmt.Printf("host: %s\n", r.Host)
	fmt.Printf("url: %s\n", r.URL.Path)
	fmt.Printf("query: %v\n", r.URL.Query())
	fmt.Printf("postform: %s\n", r.PostForm)
}

package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/config"
	"net/http"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello my world!\n")
	fmt.Printf("=> Handling request\n")
	r.ParseForm()
	fmt.Printf("method: %s\n", r.Method)
	fmt.Printf("host: %s\n", r.Host)
	fmt.Printf("url: %s\n", r.URL.Path)
	fmt.Printf("query: %v\n", r.URL.Query())
	fmt.Printf("postform: %s\n", r.PostForm)
}

func StartServer(config config.Configuration) error {
	http.HandleFunc("/", RequestHandler)
	http.ListenAndServe(fmt.Sprintf("%s:%v", config.Server.Host, config.Server.Port), nil)
	return nil
}

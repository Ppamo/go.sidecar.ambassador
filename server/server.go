package server

import (
	"fmt"
	"github.com/Ppamo/go.sidecar.ambassador/config"
	"net/http"
)

func StartServer(config config.Configuration) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello my world!\n")
	})
	http.ListenAndServe(fmt.Sprintf("%s:%v", config.Server.Host, config.Server.Port), nil)
	return nil
}

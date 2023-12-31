package web

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status": "OK"}`))
}

func NewWebServer(port int) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", Health)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	return server

}

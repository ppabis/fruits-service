package router

import (
	"context"
	"net/http"
	"strconv"
)

var server *http.Server

func Serve(port int) error {
	// Serve the web app
	mux := http.NewServeMux()
	mux.HandleFunc("/fruit", SetFruit)
	mux.HandleFunc("/login", LoginUser)
	mux.HandleFunc("/logout", LogoutUser)
	mux.HandleFunc("/token", GetToken)
	mux.HandleFunc("/metrics", PresentMetrics)
	mux.HandleFunc("/", ListAllFruits)
	server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	return server.ListenAndServe()
}

func Shutdown() {
	if server != nil {
		server.Shutdown(context.TODO())
	}
}

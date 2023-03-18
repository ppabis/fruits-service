package router

import (
	"net/http"
	"strconv"
)

func Serve(port int) error {
	// Serve the web app
	mux := http.NewServeMux()
	mux.HandleFunc("/fruit", SetFruit)
	mux.HandleFunc("/login", LoginUser)
	mux.HandleFunc("/logout", LogoutUser)
	mux.HandleFunc("/", ListAllFruits)
	return http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"urlshorten/routes"
)

func main() {
	//New router for url shortening
	r := mux.NewRouter()

	port := routes.Port
	// Routes for shortening and resolving the url
	r.HandleFunc("/shorten/", routes.ShortenUrl).Methods("POST")
	r.HandleFunc("/{id}", routes.ResolveUrl).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

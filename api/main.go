package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"urlshorten/routes"
)

func main() {
	//New router for url shortening
	r := mux.NewRouter()
	dir, _ := homedir.Dir()
	fmt.Println("dir is", dir)
	urlFilePath := filepath.Join(dir, "urls")
	err := os.MkdirAll(urlFilePath, 0744)
	if err != nil {
		log.Fatal("Failed to create directory", err)
	}
	port := routes.Port
	// Routes for shortening and resolving the url
	r.Handle("/shorten", routes.ShortenUrl(urlFilePath)).Methods("POST")
	r.Handle("/{id}", routes.ResolveUrl(urlFilePath)).Methods("GET")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}

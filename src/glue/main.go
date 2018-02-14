package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var dir string

	println("Starting glue")
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()

	// This will serve files under http://localhost:8080/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

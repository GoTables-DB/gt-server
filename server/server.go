package main

import (
	"log"
	"net/http"
	"strings"
)

const (
	PORT = ":5678"
	// CERT = "cert.pem"
	// KEY  = "key.pem"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if len(url) > 1 {
			path := strings.Split(url, "/")
			db := path[1]
			if len(path) > 2 {
				table := path[2]
			} else {
				// Print all available tables
			}
		} else {
			// Print all available databases
		}
	})
	log.Fatal(http.ListenAndServe(PORT, nil))
	// log.Fatal(http.ListenAndServeTLS(PORT, CERT, KEY, nil))
}

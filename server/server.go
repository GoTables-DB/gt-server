package main

import (
	"log"
	"net/http"
)

const (
	PORT = ":5678"
	// CERT = "cert.pem"
	// KEY  = "key.pem"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	log.Fatal(http.ListenAndServe(PORT, nil))
	// log.Fatal(http.ListenAndServeTLS(PORT, CERT, KEY, nil))
}

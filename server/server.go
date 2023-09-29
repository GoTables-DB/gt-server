package main

import (
	"encoding/json"
	"github.com/Jero075/gotables/data"
	"log"
	"net/http"
	"strings"
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
				err, tbls := data.GetTables(db)
				if err != nil {
					_, err := w.Write([]byte(err.Error()))
					log.Print(err)
				} else {
					ret, err := json.Marshal(tbls)
					if err != nil {
						_, err := w.Write([]byte(err.Error()))
						log.Print(err)
					} else {
						_, err := w.Write(ret)
						log.Print(err)
					}
				}
			}
		} else {
			err, dbs := data.GetDBs()
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				log.Print(err)
			} else {
				ret, err := json.Marshal(dbs)
				if err != nil {
					_, err := w.Write([]byte(err.Error()))
					log.Print(err)
				} else {
					_, err := w.Write(ret)
					log.Print(err)
				}
			}
		}
	})
	log.Fatal(http.ListenAndServe(config.port, nil))
	// log.Fatal(http.ListenAndServeTLS(PORT, CERT, KEY, nil))
}

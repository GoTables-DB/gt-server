package server

import (
	"encoding/json"
	"github.com/Jero075/gotables/fs"
	"log"
	"net/http"
	"strings"
)

func main() {
	config, err := fs.Config()
	if err != nil {
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if len(url) > 1 {
			path := strings.Split(url, "/")
			db := path[1]
			if len(path) > 2 {
				// table := path[2]
			} else {
				tables, tblErr := getTables(db, config.RootDir)
				if tblErr != nil {
					log.Print(tblErr)
					w.WriteHeader(500)
				} else {
					respErr := sendJson(tables, w)
					if respErr != nil {
						log.Print(respErr)
						w.WriteHeader(500)
					}
				}
			}
		} else {
			dbs, dbErr := getDBs(config.RootDir)
			if dbErr != nil {
				log.Print(dbErr)
				w.WriteHeader(500)
			} else {
				respErr := sendJson(dbs, w)
				if respErr != nil {
					log.Print(respErr)
					w.WriteHeader(500)
				}
			}
		}
	})
	if config.HTTPSMode {
		log.Fatal(http.ListenAndServeTLS(config.Port, config.SSLCert, config.SSLKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(config.Port, nil))
	}
}

func getDBs(dir string) ([]string, error) {
	dbs, dbErr := fs.GetDBs(dir)
	if dbErr != nil {
		return nil, dbErr
	} else {
		return dbs, nil
	}
}

func getTables(db, dir string) ([]string, error) {
	tables, tblErr := fs.GetTables(db, dir)
	if tblErr != nil {
		return nil, tblErr
	} else {
		return tables, nil
	}
}

func sendJson(data []string, w http.ResponseWriter) error {
	body, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return jsonErr
	}
	_, responseErr := w.Write(body)
	if responseErr != nil {
		return responseErr
	}
	return nil
}

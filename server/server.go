package server

import (
	"encoding/json"
	"github.com/Jero075/gotables/fs"
	"log"
	"net/http"
	"strings"
)

func Run(config fs.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "" {
			get(w, r, config)
		} else if r.Method == "POST" {
			post(w, r, config)
		} else if r.Method == "DELETE" {
			del(w, r, config)
		} else {
			w.WriteHeader(405)
		}
	})
	if config.HTTPSMode {
		log.Fatal(http.ListenAndServeTLS(config.Port, config.SSLCert, config.SSLKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(config.Port, nil))
	}
}

func get(w http.ResponseWriter, r *http.Request, config fs.Conf) {
	url := r.URL.Path
	if len(url) > 1 {
		path := strings.Split(url, "/")
		db := path[1]
		if len(path) > 2 {
			table := path[2]
			tbl, err, status404 := fs.GetTable(db, table, config.RootDir)
			if err != nil {
				log.Println(err)
				if status404 {
					w.WriteHeader(404)
				} else {
					w.WriteHeader(500)
				}
			} else {
				respErr := sendJson(tbl, w)
				if respErr != nil {
					log.Println(respErr)
					w.WriteHeader(500)
				}
			}
		} else {
			tables, tblErr := getTables(db, config.RootDir)
			if tblErr != nil {
				log.Println(tblErr)
				w.WriteHeader(404)
			} else {
				respErr := sendJson(tables, w)
				if respErr != nil {
					log.Println(respErr)
					w.WriteHeader(500)
				}
			}
		}
	} else {
		dbs, dbErr := getDBs(config.RootDir)
		if dbErr != nil {
			log.Println(dbErr)
			w.WriteHeader(500)
		} else {
			respErr := sendJson(dbs, w)
			if respErr != nil {
				log.Println(respErr)
				w.WriteHeader(500)
			}
		}
	}
}

func post(w http.ResponseWriter, r *http.Request, config fs.Conf) {
	w.WriteHeader(405)
}

func del(w http.ResponseWriter, r *http.Request, config fs.Conf) {
	w.WriteHeader(405)
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

func sendJson(data any, w http.ResponseWriter) error {
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

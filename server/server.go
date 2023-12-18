package server

import (
	"encoding/json"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations"
	"log"
	"net/http"
	"strconv"
)

func Run(config fs.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			table, err := operations.GTSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
			get(w, table, err)
		} else if r.Method == http.MethodHead {
			table, err := operations.GTSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
			head(w, table, err)
		} else if r.Method == http.MethodPost {
			if checkSyntaxSQL(r) {
				table, err := operations.SQLSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
				post(w, table, err)
			} else {
				table, err := operations.GTSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
				post(w, table, err)
			}
		} else if r.Method == http.MethodPut {
			table, err := operations.GTSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
			put(w, table, err)
		} else if r.Method == http.MethodDelete {
			table, err := operations.GTSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
			del(w, table, err)
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

func get(w http.ResponseWriter, table fs.Table, err error) {
	if err != nil {
		// TODO: Handle errors
		// Temporary error code
		w.WriteHeader(500)
	} else {
		jsonErr := sendTable(table, w, true)
		if jsonErr != nil {
			log.Println(jsonErr)
			w.WriteHeader(500)
		}
	}
}

func head(w http.ResponseWriter, table fs.Table, err error) {
	if err != nil {
		// TODO: Handle errors
		// Temporary error code
		w.WriteHeader(500)
	} else {
		jsonErr := sendTable(table, w, false)
		if jsonErr != nil {
			log.Println(jsonErr)
			w.WriteHeader(500)
		}
	}
}

func post(w http.ResponseWriter, table fs.Table, err error) {
	if err != nil {
		// TODO: Handle errors
		// Temporary error code
		w.WriteHeader(500)
	} else {
		jsonErr := sendTable(table, w, true)
		if jsonErr != nil {
			log.Println(jsonErr)
			w.WriteHeader(500)
		}
	}
}

func put(w http.ResponseWriter, table fs.Table, err error) {
	if err != nil {
		// TODO: Handle errors
		// Temporary error code
		w.WriteHeader(500)
	} else {
		jsonErr := sendTable(table, w, true)
		if jsonErr != nil {
			log.Println(jsonErr)
			w.WriteHeader(500)
		}
	}
}

func del(w http.ResponseWriter, table fs.Table, err error) {
	if err != nil {
		// TODO: Handle errors
		// Temporary error code
		w.WriteHeader(500)
	} else {
		jsonErr := sendTable(table, w, true)
		if jsonErr != nil {
			log.Println(jsonErr)
			w.WriteHeader(500)
		}
	}
}

func checkSyntaxSQL(r *http.Request) bool {
	if r.Header.Get("Syntax") == "SQL" {
		return true
	}
	return false
}

func sendTable(data fs.Table, w http.ResponseWriter, withBody bool) error {
	jsonData := fs.Ttoj(data)
	body, jsonErr := json.Marshal(jsonData)
	if jsonErr != nil {
		return jsonErr
	}
	if withBody {
		w.Header().Set("Content-Type", "application/json")
		_, responseErr := w.Write(body)
		if responseErr != nil {
			return responseErr
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(200)
	}
	return nil
}

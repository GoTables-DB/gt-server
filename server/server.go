package server

import (
	"encoding/json"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"log"
	"net/http"
	"strconv"
)

func Run(config fs.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var table fs.Table
		var err error
		if checkSyntaxSQL(r) {
			table, err = operations.SQLSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
		} else {
			table, err = operations.GTSyntax(r.Method, r.URL.Path, r.URL.Query().Get("query"), config)
		}
		if r.Method == http.MethodHead {
			respond(w, false, table, err)
		} else {
			respond(w, true, table, err)
		}
	})
	if config.HTTPSMode {
		log.Fatal(http.ListenAndServeTLS(config.Port, config.SSLCert, config.SSLKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(config.Port, nil))
	}
}

func respond(w http.ResponseWriter, body bool, table fs.Table, err error) {
	if err != nil {
		respondError(w, err)
	} else {
		if body {
			err = respondTable(table, w, true)
		} else {
			err = respondTable(table, w, false)
		}
		if err != nil {
			log.Println(err)
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

func respondTable(data fs.Table, w http.ResponseWriter, withBody bool) error {
	jsonData := fs.Ttoj(data)
	body, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	if withBody {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(body)
		if err != nil {
			return err
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(200)
	}
	return nil
}

func respondError(w http.ResponseWriter, err error) {
	columns := []fs.Column{{Name: "Error", Type: "str"}}
	rows := make([][]any, 0)
	rows = append(rows, []any{err.Error()})
	table, err := shared.MakeTable(columns, rows)
	if err != nil {
		w.WriteHeader(500)
	} else {
		// w.WriteHeader(err.Status)
		w.WriteHeader(500)
		tbl, err := json.Marshal(fs.Ttoj(table))
		if err != nil {
			w.WriteHeader(500)
		}
		_, err = w.Write(tbl)
		if err != nil {
			w.WriteHeader(500)
		}
	}
}

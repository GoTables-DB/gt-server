package server

import (
	"encoding/json"
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Body struct {
	Query     string `json:"query"`
	SessionId string `json:"session_id"`
}

func Run(config fs.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var tbl table.Table
		var err error
		q, err := query(r)
		if err != nil && r.Method == http.MethodPost {
			respondError(w, errors.New("failed to read request Body: "+err.Error()))
		}
		if checkSyntaxSQL(r) {
			tbl, err = operations.SQLSyntax(r.Method, r.URL.Path, q, config)
		} else {
			tbl, err = operations.GTSyntax(r.Method, r.URL.Path, q, config)
		}
		if r.Method == http.MethodHead {
			respond(w, false, tbl, err)
		} else {
			respond(w, true, tbl, err)
		}
	})
	if config.HTTPSMode {
		log.Fatal(http.ListenAndServeTLS(config.Port, config.SSLCert, config.SSLKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(config.Port, nil))
	}
}

func respond(w http.ResponseWriter, body bool, table table.Table, err error) {
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

func respondTable(data table.Table, w http.ResponseWriter, withBody bool) error {
	jsonData := data.ToU()
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
	columns := []table.Column{{Name: "Error", Type: "str"}}
	rows := make([]map[string]any, 0)
	row := map[string]any{}
	row["Error"] = err.Error()
	rows = append(rows, row)
	tbl, err := shared.MakeTable(columns, rows)
	if err != nil {
		w.WriteHeader(500)
	} else {
		// w.WriteHeader(err.Status)
		w.WriteHeader(404)
		tbl, err := json.Marshal(tbl.ToU())
		if err != nil {
			w.WriteHeader(500)
		}
		_, err = w.Write(tbl)
		if err != nil {
			w.WriteHeader(500)
		}
	}
}

func query(r *http.Request) (string, error) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	queryStruct := Body{}
	err = json.Unmarshal(data, &queryStruct)
	if err != nil {
		return "", err
	}
	return queryStruct.Query, nil
}

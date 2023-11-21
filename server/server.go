package server

import (
	"encoding/json"
	"github.com/Jero075/gotables/fs"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PostBody struct {
	Name string `json:"name"`
}

type PutBody struct {
	Name        string `json:"name"`
	CellsPerRow int    `json:"row_length"`
}

func Run(config fs.Conf) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			get(w, r, config, true)
		} else if r.Method == http.MethodHead {
			get(w, r, config, false)
		} else if r.Method == http.MethodPost {
			post(w, r, config)
		} else if r.Method == http.MethodPut {
			put(w, r, config)
		} else if r.Method == http.MethodDelete {
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

func get(w http.ResponseWriter, r *http.Request, config fs.Conf, withBody bool) {
	db, table := url(r)
	if db == "" {
		dbs, dbErr := getDBs(config.Dir)
		if dbErr != nil {
			log.Println(dbErr)
			w.WriteHeader(500)
		} else {
			respErr := sendJson(dbs, w, withBody)
			if respErr != nil {
				log.Println(respErr)
				w.WriteHeader(500)
			}
		}
	} else if table == "" {
		tables, tblErr := getTables(db, config.Dir)
		if tblErr != nil {
			log.Println(tblErr)
			w.WriteHeader(404)
		} else {
			respErr := sendJson(tables, w, withBody)
			if respErr != nil {
				log.Println(respErr)
				w.WriteHeader(500)
			}
		}
	} else {
		tbl, err, status404 := fs.GetTable(db, table, config.Dir)
		if err != nil {
			log.Println(err)
			if status404 {
				w.WriteHeader(404)
			} else {
				w.WriteHeader(500)
			}
		} else {
			respErr := sendJson(tbl, w, withBody)
			if respErr != nil {
				log.Println(respErr)
				w.WriteHeader(500)
			}
		}
	}
}

func post(w http.ResponseWriter, r *http.Request, config fs.Conf) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(415)
		return
	}
	db, table := url(r)
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) < 2 {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	postBody := PostBody{}
	jsonErr := json.Unmarshal(body, &postBody)
	if jsonErr != nil {
		log.Println(jsonErr)
		w.WriteHeader(500)
		return
	}
	if db == "" || table == "" {
		w.WriteHeader(400)
	} else {
		// TODO: Possibility to append content to table
	}
}

func put(w http.ResponseWriter, r *http.Request, config fs.Conf) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(415)
		return
	}
	db, table := url(r)
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) < 2 {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	putBody := PutBody{}
	jsonErr := json.Unmarshal(body, &putBody)
	if jsonErr != nil {
		log.Println(jsonErr)
		w.WriteHeader(500)
		return
	}
	if db == "" {
		if putBody.Name == "" {
			w.WriteHeader(400)
			return
		}
		fsErr := fs.NewDB(putBody.Name, config.Dir)
		if fsErr != nil {
			log.Println(fsErr)
			w.WriteHeader(500)
			return
		}
	} else if table == "" {
		if putBody.Name == "" {
			w.WriteHeader(400)
			return
		}
		dir := config.Dir + "/" + db
		fsErr := fs.NewTable(putBody.Name, dir, putBody.CellsPerRow)
		if fsErr != nil {
			log.Println(fsErr)
			w.WriteHeader(500)
			return
		}
	} else {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(201)
}

func del(w http.ResponseWriter, r *http.Request, config fs.Conf) {
	w.WriteHeader(405)
}

func url(r *http.Request) (db string, table string) {
	path := r.URL.EscapedPath()
	splitPath := strings.Split(path, "/")
	var splitPathNoEmpty []string
	for _, element := range splitPath {
		if element != "" {
			splitPathNoEmpty = append(splitPathNoEmpty, element)
		}
	}
	if len(splitPathNoEmpty) > 0 {
		db = splitPathNoEmpty[0]
		if len(splitPathNoEmpty) > 1 {
			table = splitPathNoEmpty[1]
		} else {
			table = ""
		}
	} else {
		db = ""
		table = ""
	}
	return db, table
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

func sendJson(data any, w http.ResponseWriter, withBody bool) error {
	body, jsonErr := json.Marshal(data)
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

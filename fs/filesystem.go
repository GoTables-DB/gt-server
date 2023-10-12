package fs

import (
	"crypto/sha512"
	"encoding/json"
	"os"
)

type DB struct {
	name    string
	tables  []Table
	dbusers map[string][64]byte
}

type Table struct {
	name    string
	access  map[string]int
	indexes []map[string]int
	table   []map[string]any
}

type Conf struct {
	Port      string `json:"port"`
	RootDir   string `json:"root_dir"`
	HTTPSMode bool   `json:"https_mode"`
	SSLCert   string `json:"ssl_cert"`
	SSLKey    string `json:"ssl_key"`
}

func NewDB(name string, adminUsername string, adminPassword string) error {
	pwHash := sha512.Sum512([]byte(adminPassword))
	db := DB{name: name, tables: make([]Table, 0), dbusers: make(map[string][64]byte)}
	db.dbusers[adminUsername] = pwHash

	return nil
}

func GetDBs(dir string) ([]string, error) {
	dbs, err := ls(dir)
	return dbs, err
}

func GetTables(db, dir string) ([]string, error) {
	tables, err := ls(dir + "/" + db)
	return tables, err
}

func GetTable(db, table, dir string) (Table, error, bool) {
	tableFile, err := os.ReadFile(dir + "/" + db + "/" + table + ".json")
	if err != nil {
		return Table{}, err, true
	}
	tableData := Table{}
	jsonErr := json.Unmarshal(tableFile, &tableData)
	if jsonErr != nil {
		return Table{}, jsonErr, false
	}
	return tableData, nil, false
}

func Config() (Conf, error) {
	confFile, fileErr := os.ReadFile("gtconfig.json")
	if fileErr != nil {
		return Conf{}, fileErr
	}
	config := Conf{}
	jsonErr := json.Unmarshal(confFile, &config)
	if jsonErr != nil {
		return Conf{}, jsonErr
	}
	return config, nil
}

func ls(dir string) (contents []string, error error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	contents = []string{}
	for _, entry := range entries {
		contents = append(contents, entry.Name())
	}
	return contents, nil
}

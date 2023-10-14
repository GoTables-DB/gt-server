package fs

import (
	"encoding/json"
	"os"
	"strings"
)

type DB struct {
	dbUsers map[string][64]byte
	access  []map[string]int
}

type Table struct {
	Rows     []map[string]any `json:"rows"`
	Defaults map[string]any   `json:"defaults"`
	// indexes  []map[string]int `json:"indexes"`
}

type Conf struct {
	Port      string `json:"port"`
	RootDir   string `json:"root_dir"`
	HTTPSMode bool   `json:"https_mode"`
	SSLCert   string `json:"ssl_cert"`
	SSLKey    string `json:"ssl_key"`
}

/*
func NewDB() error {
	return nil
}

func NewTable() error {
	return nil
}
*/

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
		contents = append(contents, strings.TrimSuffix(entry.Name(), ".json"))
	}
	return contents, nil
}

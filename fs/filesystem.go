package fs

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Table struct {
	Rows     [][]any        `json:"rows"`
	Defaults map[string]any `json:"defaults"`
}

type Conf struct {
	Port            string `json:"port"`
	Dir             string `json:"dir"`
	HTTPSMode       bool   `json:"https_mode"`
	SSLCert         string `json:"ssl_cert"`
	SSLKey          string `json:"ssl_key"`
	EnableGTSyntax  bool   `json:"enable_gt_syntax"`
	EnableSQLSyntax bool   `json:"enable_sql_syntax"`
	EnableGoSyntax  bool   `json:"enable_go_syntax"`
}

func NewDB(name string, dir string) error {
	dbLocation := dir + "/" + name
	err := os.Mkdir(dbLocation, 0755)
	if err != nil {
		return err
	}
	return nil
}

func NewTable(name string, dir string, rowLen int) error {
	tblLocation := dir + "/" + name + ".json"
	tbl := Table{Rows: make([][]any, rowLen)}
	data, jsonErr := json.Marshal(tbl)
	if jsonErr != nil {
		return jsonErr
	}
	fsErr := os.WriteFile(tblLocation, data, 0755)
	if fsErr != nil {
		return fsErr
	}
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

func GetTable(db, table, dir string) (Table, error) {
	if table == "" {
		return Table{}, errors.New("no table specified")
	}
	tableFile, err := os.ReadFile(dir + "/" + db + "/" + table + ".json")
	if err != nil {
		return Table{}, errors.New("table not found")
	}
	tableData := Table{}
	jsonErr := json.Unmarshal(tableFile, &tableData)
	if jsonErr != nil {
		return Table{}, jsonErr
	}
	return tableData, nil
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

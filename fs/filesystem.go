package fs

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

type Column struct {
	Name string `json:"name"`
	Type any    `json:"type"` // Any value of a specific datatype. reflect.TypeOf() to gt-get the type.
}

type Table struct {
	ColumnNames []Column        `json:"column_names"`
	Rows        [][]interface{} `json:"rows"` // Row 1 for defaults
}

type Conf struct {
	// Basic config
	Port   string `json:"port"`
	Dir    string `json:"dir"`
	LogDir string `json:"log_dir"`
	// HTTPS config
	HTTPSMode bool   `json:"https"`
	SSLCert   string `json:"cert"`
	SSLKey    string `json:"key"`
	// Query config
	EnableGTSyntax  bool `json:"gt_syntax"`
	EnableSQLSyntax bool `json:"sql_syntax"`
	// Advanced config
	// ConnectionTimeout int 'json:"conn_timeout"`
	// MaxConnections int `json:"conn_max"`
}

func NewDB(name string, dir string) error {
	dbLocation := dir + "/" + name
	err := os.Mkdir(dbLocation, 0755)
	if err != nil {
		return err
	}
	return nil
}

func NewTable(name string, dir string) error {
	tblLocation := dir + "/" + name + ".json"
	tbl := Table{}
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
	// Defaults
	config := Conf{
		Port:            ":5678",
		Dir:             "/srv/GoTables/server",
		LogDir:          "/srv/GoTables/logs",
		EnableGTSyntax:  true,
		EnableSQLSyntax: true,
	}
	if _, err := os.Stat("gtconfig.json"); err == nil {
		confFile, fileErr := os.ReadFile("gtconfig.json")
		if fileErr != nil {
			return Conf{}, fileErr
		}
		jsonErr := json.Unmarshal(confFile, &config)
		if jsonErr != nil {
			return Conf{}, jsonErr
		}
	} else {
		log.Println("Warning: configuration file not found. Using default config.")
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

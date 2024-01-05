package fs

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Column struct {
	Name    string `json:"name"`
	Type    any    `json:"type"` // Any value of a specific datatype. reflect.TypeOf() to gt-get the type.
	Default any    `json:"default"`
}

type TableJSON struct {
	Columns []Column        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type Table struct {
	columns []Column
	rows    [][]interface{} // Row 1 for defaults
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

func DetermineDatatype(datatype string) (any, error) {
	var ret any
	var err error
	switch datatype {
	// String
	case "str":
		ret = ""
	// Integer
	case "int":
		ret = 0
	// Float
	case "flt":
		ret = 0.0
	// Boolean
	case "bol":
		ret = false
	// Date
	case "dat":
		ret = time.Time{}
	// Table
	case "tab":
		ret = Table{}
	default:
		err = errors.New("unknown datatype")
	}
	return ret, err
}

/// Methods for Table ///

func (t Table) GetColumns() []Column {
	return t.columns
}

func (t Table) GetRows() [][]interface{} {
	return t.rows
}

func (t Table) SetColumns(columns []Column) Table {
	t.columns = columns
	return t
}

func (t Table) SetRows(rows [][]interface{}) (Table, error) {
	for i, row := range rows {
		if len(row) != len(t.columns) {
			return Table{}, errors.New("row length of row " + strconv.Itoa(i) + " is invalid")
		}
		for j, cell := range row {
			if reflect.TypeOf(cell) != reflect.TypeOf(t.columns[j].Type) {
				return Table{}, errors.New("type of cell " + strconv.Itoa(j) + " in row " + strconv.Itoa(i) + " is invalid")
			}
		}
	}
	t.rows = rows
	return t, nil
}

/// Convert between TableJSON and Table ///

// Jtot - JSON to Table
func Jtot(j TableJSON) (Table, error) {
	t := Table{}
	t = t.SetColumns(j.Columns)
	t, err := t.SetRows(j.Rows)
	return t, err
}

// Ttoj - Table to JSON
func Ttoj(t Table) TableJSON {
	j := TableJSON{}
	j.Columns = t.GetColumns()
	j.Rows = t.GetRows()
	return j
}

/// Read and write to filesystem ///

func NewDB(name, dir string) error {
	dbLocation := dir + "/" + name
	err := os.Mkdir(dbLocation, 0755)
	return err
}

func NewTable(name, db, dir string) error {
	err := writeTable(Table{}, name, db, dir)
	return err
}

func DeleteDB(name, dir string) error {
	dbLocation := dir + "/" + name
	err := os.RemoveAll(dbLocation)
	return err
}

func DeleteTable(name, db, dir string) error {
	tblLocation := dir + "/" + db + "/" + name + ".json"
	err := os.Remove(tblLocation)
	return err
}

func GetDBs(dir string) ([]string, error) {
	dbs, err := ls(dir)
	return dbs, err
}

func GetTables(db, dir string) ([]string, error) {
	tables, err := ls(dir + "/" + db)
	return tables, err
}

func GetTable(name, db, dir string) (Table, error) {
	if name == "" {
		return Table{}, errors.New("no table specified")
	}
	tableFile, err := os.ReadFile(dir + "/" + db + "/" + name + ".json")
	if err != nil {
		return Table{}, errors.New("table " + name + " in database " + db + " could not be found")
	}
	tableData := TableJSON{}
	err = json.Unmarshal(tableFile, &tableData)
	if err != nil {
		return Table{}, err
	}
	table, err := Jtot(tableData)
	return table, err
}

func MoveDB(oldName, name, dir string) error {
	err := CopyDB(oldName, name, dir)
	if err != nil {
		return err
	}
	err = DeleteDB(oldName, dir)
	return err
}

func CopyDB(oldName, name, dir string) error {
	exists, err := existsDB(name, dir)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("database " + name + " already exists")
	}
	err = NewDB(name, dir)
	if err != nil {
		return err
	}
	err = cpDB(oldName, name, dir)
	return err
}

func MoveTable(oldName, name, db, dir string) error {
	err := CopyTable(oldName, name, db, dir)
	if err != nil {
		return err
	}
	err = DeleteTable(oldName, db, dir)
	return err
}

func CopyTable(oldName, name, db, dir string) error {
	exists, err := existsTable(name, db, dir)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("table" + name + " in database " + db + " already exists")
	}
	err = NewTable(name, db, dir)
	if err != nil {
		return err
	}
	err = cpTable(oldName, name, db, dir)
	return err
}

func ModifyTable(data Table, name string, db string, dir string) error {
	_, err := GetTable(name, db, dir)
	if err != nil {
		return err
	}
	err = writeTable(data, name, db, dir)
	return err
}

/// Load config ///

func Config() (Conf, error) {
	// Defaults
	config := Conf{
		Port:            ":5678",
		Dir:             "/srv/GoTables/server",
		LogDir:          "/srv/GoTables/logs",
		EnableGTSyntax:  true,
		EnableSQLSyntax: true,
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Conf{}, err
	}
	if _, err := os.Stat(home + "/.config/gotables/config.json"); err == nil {
		confFile, fileErr := os.ReadFile(home + "/.config/gotables/config.json")
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

/// Helper functions ///

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

func existsDB(name, dir string) (bool, error) {
	dbs, err := GetDBs(dir)
	if err != nil {
		return false, err
	}
	for _, db := range dbs {
		if db == name {
			return true, nil
		}
	}
	return false, nil
}

func cpDB(oldName, name, dir string) error {
	files, err := ls(dir + "/" + oldName)
	if err != nil {
		return err
	}
	for i := 0; i < len(files); i++ {
		file, err := os.ReadFile(dir + "/" + oldName + "/" + files[i] + ".json")
		if err != nil {
			return err
		}
		err = os.WriteFile(dir+"/"+name+"/"+files[i]+".json", file, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func existsTable(name, db, dir string) (bool, error) {
	tables, err := GetTables(db, dir)
	if err != nil {
		return false, err
	}
	for _, table := range tables {
		if table == name {
			return true, nil
		}
	}
	return false, nil
}

func cpTable(oldName, name, db, dir string) error {
	file, err := os.ReadFile(dir + "/" + db + "/" + oldName + ".json")
	if err != nil {
		return err
	}
	err = os.WriteFile(dir+"/"+db+"/"+name+".json", file, 0755)
	return err
}

func writeTable(data Table, name string, db string, dir string) error {
	tblLocation := dir + "/" + db + "/" + name + ".json"
	tbl := TableJSON{
		Columns: data.columns,
		Rows:    data.rows,
	}
	jsonData, jsonErr := json.Marshal(tbl)
	if jsonErr != nil {
		return jsonErr
	}
	err := os.WriteFile(tblLocation, jsonData, 0755)
	return err
}

package fs

import (
	"encoding/json"
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"log"
	"os"
	"strings"
)

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

func AddDB(name, dir string) error {
	dbLocation := dir + "/" + name
	err := os.Mkdir(dbLocation, 0755)
	return err
}

func AddTable(name, db, dir string) error {
	err := writeTable(table.Table{}, name, db, dir)
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

func GetTable(name, db, dir string) (table.Table, error) {
	if name == "" {
		return table.Table{}, errors.New("no table specified")
	}
	tableFile, err := os.ReadFile(dir + "/" + db + "/" + name + ".json")
	if err != nil {
		return table.Table{}, errors.New("table " + name + " in database " + db + " could not be found")
	}
	tableData := table.TableU{}
	err = json.Unmarshal(tableFile, &tableData)
	if err != nil {
		return table.Table{}, err
	}
	tbl, err := tableData.ToT()
	return tbl, err
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
	err = AddDB(name, dir)
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
	err = AddTable(name, db, dir)
	if err != nil {
		return err
	}
	err = cpTable(oldName, name, db, dir)
	return err
}

func ModifyTable(data table.Table, name string, db string, dir string) error {
	_, err := GetTable(name, db, dir)
	if err != nil {
		return err
	}
	err = writeTable(data, name, db, dir)
	return err
}

/// Load config ///

func Config(location string) (Conf, error) {
	// Defaults
	config := Conf{
		Port:            ":5678",
		Dir:             "/srv/GoTables/server",
		LogDir:          "/srv/GoTables/logs",
		EnableGTSyntax:  true,
		EnableSQLSyntax: true,
	}
	if location == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			location = home + "/.config/gotables/config.json"
		}
	}
	if _, err := os.Stat(location); err == nil {
		confFile, fileErr := os.ReadFile(location)
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
	for _, tbl := range tables {
		if tbl == name {
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

func writeTable(data table.Table, name string, db string, dir string) error {
	tblLocation := dir + "/" + db + "/" + name + ".json"
	tbl := table.TableU{
		Columns: data.GetColumns(),
		Rows:    data.GetRows(),
	}
	jsonData, jsonErr := json.Marshal(tbl)
	if jsonErr != nil {
		return jsonErr
	}
	err := os.WriteFile(tblLocation, jsonData, 0755)
	return err
}

package gt_get

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
)

func Get(table string, db string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error = nil

	if db == "" {
		retTable, retError = getDBs(config.Dir)
	} else if table == "" {
		retTable, retError = getTables(db, config.Dir)
	} else {
		retTable, retError = fs.GetTable(table, db, config.Dir)
	}

	return retTable, retError
}

func getDBs(dir string) (fs.Table, error) {
	dbs, err := fs.GetDBs(dir)
	if err != nil {
		return fs.Table{}, err
	}
	columns := []fs.Column{{Name: "Databases", Type: "string"}}
	rows := make([][]interface{}, 0)
	for _, db := range dbs {
		rows = append(rows, []interface{}{db})
	}
	return shared.MakeNewTable(columns, rows)
}

func getTables(db string, dir string) (fs.Table, error) {
	tables, err := fs.GetTables(db, dir)
	if err != nil {
		return fs.Table{}, err
	}
	columns := []fs.Column{{Name: "Tables", Type: "string"}}
	rows := make([][]interface{}, 0)
	for _, table := range tables {
		rows = append(rows, []interface{}{table})
	}
	return shared.MakeNewTable(columns, rows)
}

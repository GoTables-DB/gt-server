package gt_get

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations"
)

func Get(db string, table string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error = nil

	if db == "" {

	} else if table == "" {
		retTable, retError = getTables(db, config.Dir)
	} else {

	}

	return retTable, retError
}

func getTables(db string, dir string) (fs.Table, error) {
	tables, err := fs.GetTables(db, dir)
	if err != nil {
		return fs.Table{}, err
	}
	column := fs.Column{
		Name: "Tables",
		Type: fs.Table{},
	}
	columns := []fs.Column{column}
	rows := make([][]interface{}, 0)
	for i, table := range tables {
		rows[i] = append(rows[i], table)
	}
	retTable := operations.MakeTableNew(columns, rows)
	return retTable, nil
}

package gt_get

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"git.jereileu.ch/gotables/server/gt-server/table"
)

func Get(tbl string, db string, config fs.Conf) (table.Table, error) {
	retTable := table.Table{}
	var retError error = nil

	if db == "" {
		retTable, retError = getDBs(config.Dir)
	} else if tbl == "" {
		retTable, retError = getTables(db, config.Dir)
	} else {
		retTable, retError = fs.GetTable(tbl, db, config.Dir)
	}

	return retTable, retError
}

func getDBs(dir string) (table.Table, error) {
	dbs, err := fs.GetDBs(dir)
	if err != nil {
		return table.Table{}, err
	}
	columns := []table.Column{{Name: "Databases", Type: "str"}}
	rows := make([]map[string]any, 0)
	for _, db := range dbs {
		row := map[string]any{}
		row["Databases"] = db
		rows = append(rows, row)
	}
	return shared.MakeTable(columns, rows)
}

func getTables(db string, dir string) (table.Table, error) {
	tables, err := fs.GetTables(db, dir)
	if err != nil {
		return table.Table{}, err
	}
	columns := []table.Column{{Name: "Tables", Type: "str"}}
	rows := make([]map[string]any, 0)
	for _, tbl := range tables {
		row := map[string]any{}
		row["Tables"] = tbl
		rows = append(rows, row)
	}
	return shared.MakeTable(columns, rows)
}

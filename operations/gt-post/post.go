package gt_post

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"strings"
)

func Post(query []string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error

	switch strings.ToLower(query[0]) {
	// List dbs
	case "list":
		if len(query) != 1 {
			retError = errors.New("invalid syntax")
		} else {
			dbs, err := fs.GetDBs(config.Dir)
			if err != nil {
				return fs.Table{}, err
			}
			columns := []fs.Column{{Name: "Databases", Type: "string"}}
			rows := make([][]interface{}, 0)
			for _, db := range dbs {
				rows = append(rows, []interface{}{db})
			}
			retTable, retError = shared.MakeTable(columns, rows)
		}
	case "database":
		if len(query) < 3 {
			retError = errors.New("invalid syntax")
		}
		switch strings.ToLower(query[2]) {
		case "create":
			if len(query) != 3 {
				retError = errors.New("invalid syntax")
			} else {
				retError = fs.NewDB(query[1], config.Dir)
			}
		case "list":
			if len(query) != 3 {
				retError = errors.New("invalid syntax")
			} else {
				tables, err := fs.GetTables(query[1], config.Dir)
				if err != nil {
					return fs.Table{}, err
				}
				columns := []fs.Column{{Name: "Tables", Type: "string"}}
				rows := make([][]interface{}, 0)
				for _, table := range tables {
					rows = append(rows, []interface{}{table})
				}
				retTable, retError = shared.MakeTable(columns, rows)
			}
		case "delete":
			if len(query) != 3 {
				retError = errors.New("invalid syntax")
			}
			retError = fs.DeleteDB(query[1], config.Dir)
		case "table":
			if len(query) < 5 {
				retError = errors.New("invalid syntax")
			}
			switch strings.ToLower(query[4]) {
			case "create":
				if len(query) != 5 {
					if len(query) < 7 {
						retError = errors.New("invalid syntax")
					} else {
						if query[5] != "columns" {
							retError = errors.New("invalid syntax")
						} else {
							retTable, retError = makeTableWithColumns(query[6:], query[3], query[1], config.Dir)
						}
					}
				} else {
					retError = fs.NewTable(query[3], query[1], config.Dir)
				}
			case "show":
			case "modify":
			case "delete":
				if len(query) != 5 {
					retError = errors.New("invalid syntax")
				} else {
					retError = fs.DeleteTable(query[3], query[1], config.Dir)
				}
			case "column":
			}
		}
	case "user":
		// TODO: Implement user management
	default:
		retError = errors.New("invalid syntax")
	}

	return retTable, retError
}

func makeTableWithColumns(columns []string, table string, db string, dir string) (fs.Table, error) {
	err := fs.NewTable(table, db, dir)
	if err != nil {
		return fs.Table{}, err
	}
	tbl, err := shared.MakeTableWithColumns(columns)
	if err != nil {
		return fs.Table{}, err
	}
	err = fs.ModifyTable(tbl, table, db, dir)
	return tbl, err
}

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
	case "show":
		if len(query) != 1 {
			return fs.Table{}, errors.New("invalid syntax")
		}
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
	case "database":
		if len(query) < 3 {
			return fs.Table{}, errors.New("invalid syntax")
		}
		switch strings.ToLower(query[2]) {
		case "create":
			if len(query) != 3 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.NewDB(query[1], config.Dir)
		case "show":
			if len(query) != 3 {
				return fs.Table{}, errors.New("invalid syntax")
			}
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
		case "move":
			if len(query) != 4 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.MoveDB(query[1], query[3], config.Dir)
		case "copy":
			if len(query) != 4 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.CopyDB(query[1], query[3], config.Dir)
		case "delete":
			if len(query) != 3 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.DeleteDB(query[1], config.Dir)
		case "table":
			if len(query) < 5 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			switch strings.ToLower(query[4]) {
			case "create":
				if len(query) != 5 {
					if len(query) < 7 {
						return fs.Table{}, errors.New("invalid syntax")
					}
					if query[5] != "columns" {
						return fs.Table{}, errors.New("invalid syntax")
					}
					retTable, retError = makeTableWithColumns(query[6:], query[3], query[1], config.Dir)
				} else {
					retError = fs.NewTable(query[3], query[1], config.Dir)
				}
			case "show":
				table, err := fs.GetTable(query[3], query[1], config.Dir)
				if err != nil {
					return fs.Table{}, err
				}
				if len(query) == 5 {
					retTable = table
				} else if len(query) == 6 {
					retTable, retError = showTable(query[5], table, []string{})
				} else if query[6] == "condition" {
					return fs.Table{}, errors.New("not yet implemented")
					// retTable, retError = showTable(query[5], table, query[7])
				} else {
					return fs.Table{}, errors.New("invalid syntax")
				}
			case "move":
				if len(query) != 6 {
					return fs.Table{}, errors.New("invalid syntax")
				}
				retError = fs.MoveTable(query[3], query[5], query[1], config.Dir)
			case "copy":
				if len(query) != 6 {
					return fs.Table{}, errors.New("invalid syntax")
				}
				retError = fs.CopyTable(query[3], query[5], query[1], config.Dir)
			case "delete":
				if len(query) != 5 {
					return fs.Table{}, errors.New("invalid syntax")
				}
				retError = fs.DeleteTable(query[3], query[1], config.Dir)
			case "column":
				if len(query) < 7 {
					return fs.Table{}, errors.New("invalid syntax")
				}
			case "row":
			}
		}
	case "user":
		// TODO: Implement user management
	case "backup":
		// TODO: Implement backups
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

func showTable(columns string, table fs.Table, condition []string) (fs.Table, error) {
	colSlice := strings.Split(columns, ":")
	rows := make([]int, 0)
	for i := 0; i < len(table.GetRows()); i++ {
		rows = append(rows, i)
	}
	cols, err := shared.SelectColumns(colSlice, table)
	if err != nil {
		return fs.Table{}, err
	}
	retTable, err := shared.MakeTableFromTable(cols, rows, table)
	return retTable, err
}

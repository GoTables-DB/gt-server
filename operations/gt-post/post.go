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
				if len(query) == 5 { // Show entire table
					retTable = table
				} else if len(query) == 6 { // Show specific columns
					retTable, retError = showTable(query[5], table, []string{})
				} else if query[6] == "where" { // Show specific columns (with condition)
					retTable, retError = showTable(query[5], table, query[7:])
				} else { // Invalid syntax
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
	rows := make([]int, 0)
	for i := 0; i < len(table.GetRows()); i++ {
		rows = append(rows, i)
	}
	if len(condition) != 0 {
		var indices []int
		if (len(condition)+1)%4 != 0 {
			return fs.Table{}, errors.New("invalid condition")
		}
		for i := 0; i < len(table.GetRows()); i++ {
			var results []bool
			for j := 0; j < (len(condition)+1)/4; j++ {
				result, err := checkCondition(i, table, condition[i*4:3+i*4])
				if err != nil {
					return fs.Table{}, err
				}
				results = append(results, result)
			}
			if checkResults(results, condition) {
				indices = append(indices, i)
			}
		}
		rows = indices
	}
	colSlice := strings.Split(columns, ":")
	cols, err := shared.SelectColumns(colSlice, table)
	if err != nil {
		return fs.Table{}, err
	}
	retTable, err := shared.MakeTableFromTable(cols, rows, table)
	return retTable, err
}

func checkCondition(row int, table fs.Table, condition []string) (bool, error) {
	if len(condition) != 3 {
		return false, errors.New("invalid condition")
	}
	cols := table.GetColumns()
	rows := table.GetRows()
	isString := []bool{
		strings.HasPrefix(condition[0], "\"") && strings.HasSuffix(condition[0], "\"") || strings.HasPrefix(condition[0], "'") && strings.HasSuffix(condition[0], "'"),
		strings.HasPrefix(condition[2], "\"") && strings.HasSuffix(condition[2], "\"") || strings.HasPrefix(condition[2], "'") && strings.HasSuffix(condition[2], "'"),
	}
	var values []string
	if !isString[0] {
		col := getColumnIndex(condition[0], cols)
		if col == -1 {
			return false, errors.New("invalid condition: column " + condition[0] + " does not exist")
		}
		values = append(values, rows[row][col].(string))
	} else {
		values = append(values, condition[0])
	}
	if !isString[1] {
		col := getColumnIndex(condition[2], cols)
		if col == -1 {
			return false, errors.New("invalid condition: column " + condition[2] + " does not exist")
		}
		values = append(values, rows[row][col].(string))
	} else {
		values = append(values, condition[2])
	}
	switch condition[1] {
	case "==":
		if values[0] == values[1] {
			return true, nil
		}
		return false, nil
	case "!=":
		if values[0] != values[1] {
			return true, nil
		}
		return false, nil
	case "<":
		if values[0] < values[1] {
			return true, nil
		}
		return false, nil
	case ">":
		if values[0] > values[1] {
			return true, nil
		}
		return false, nil
	case "<=":
		if values[0] <= values[1] {
			return true, nil
		}
		return false, nil
	case ">=":
		if values[0] >= values[1] {
			return true, nil
		}
		return false, nil
	default:
		return false, errors.New("invalid condition")
	}
}

func checkResults(results []bool, condition []string) bool {
	if len(results) == 1 {
		return results[0]
	}
	var operators []string
	for i := 3; i < len(condition); i += 4 {
		operators = append(operators, condition[i])
	}
	var res [][]bool
	var currentRes []bool
	for i := 0; i < len(results); i++ {
		currentRes = append(currentRes, results[i])
		if i == len(results)-1 {
			res = append(res, currentRes)
		} else if operators[i] == "||" {
			res = append(res, currentRes)
			currentRes = make([]bool, 0)
		}
	}
	var resultsFinal []bool
	for _, result := range res {
		for i := 0; i < len(result); i++ {
			if !result[i] {
				resultsFinal = append(resultsFinal, false)
			}
		}
	}
	for _, resultFinal := range resultsFinal {
		if resultFinal {
			return true
		}
	}
	return false
}

func getColumnIndex(name string, cols []fs.Column) int {
	for i := 0; i < len(cols); i++ {
		if name == cols[i].Name {
			return i
		}
	}
	return -1
}

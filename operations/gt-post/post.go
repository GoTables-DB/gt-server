package gt_post

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"strings"
)

func Post(query []string, table string, db string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error

	if len(query) < 1 {
		return fs.Table{}, errors.New("invalid syntax")
	}
	if db == "" {
		switch strings.ToLower(query[0]) {
		case "show":
			if len(query) != 1 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			dbs, err := fs.GetDBs(config.Dir)
			if err != nil {
				return fs.Table{}, err
			}
			retTable, retError = simpleTable("Databases", dbs)
		case "user":
			if len(query) != 0 {
				return fs.Table{}, errors.New("invalid syntax")
			}
		case "backup":
			if len(query) != 0 {
				return fs.Table{}, errors.New("invalid syntax")
			}
		default:
			retError = errors.New("invalid syntax")
		}
	} else if table == "" {
		switch strings.ToLower(query[0]) {
		case "show":
			if len(query) != 1 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			tables, err := fs.GetTables(db, config.Dir)
			if err != nil {
				return fs.Table{}, err
			}
			retTable, retError = simpleTable("Tables", tables)
		case "create":
			if len(query) != 1 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.NewDB(db, config.Dir)
		case "move":
			if len(query) != 2 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.MoveDB(db, query[1], config.Dir)
		case "copy":
			if len(query) != 2 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.CopyDB(db, query[1], config.Dir)
		case "delete":
			if len(query) != 1 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.DeleteDB(db, config.Dir)
		default:
			retError = errors.New("invalid syntax")
		}
	} else {
		switch strings.ToLower(query[0]) {
		case "show":
			tbl, err := fs.GetTable(table, db, config.Dir)
			if err != nil {
				return fs.Table{}, err
			}
			if len(query) == 1 { // Show entire table
				retTable = tbl
			} else if len(query) == 2 { // Show specific columns
				retTable, retError = showTable(query[1], tbl, []string{})
			} else if query[2] == "where" { // Show specific columns (with condition)
				retTable, retError = showTable(query[1], tbl, query[3:])
			} else { // Invalid syntax
				return fs.Table{}, errors.New("invalid syntax")
			}
		case "create":
			if len(query) > 1 {
				retTable, retError = makeTableWithColumns(query[1:], table, db, config.Dir)
			} else {
				retError = fs.NewTable(table, db, config.Dir)
			}
		case "move":
			if len(query) != 2 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.MoveTable(table, query[1], db, config.Dir)
		case "copy":
			if len(query) != 2 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.CopyTable(table, query[1], db, config.Dir)
		case "delete":
			if len(query) != 1 {
				return fs.Table{}, errors.New("invalid syntax")
			}
			retError = fs.DeleteTable(table, db, config.Dir)
		case "column":
		case "row":
		default:
			retError = errors.New("invalid syntax")
		}
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

// Used to display names of databases or names of tables in a db
func simpleTable(colName string, rows []string) (fs.Table, error) {
	columns := []fs.Column{{Name: colName, Type: "str"}}
	rowSlice := make([][]interface{}, 0)
	for _, row := range rows {
		rowSlice = append(rowSlice, []interface{}{row})
	}
	return shared.MakeTable(columns, rowSlice)
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
				result, err := checkCondition(i, table, condition[j*4:3+j*4])
				if err != nil {
					return fs.Table{}, err
				}
				results = append(results, result)
			}
			var operators []string
			for j := 3; j < len(condition); j += 4 {
				operators = append(operators, condition[j])
			}
			if checkResults(results, operators) {
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
	isVar := []bool{
		strings.HasPrefix(condition[0], "\"") && strings.HasSuffix(condition[0], "\"") || strings.HasPrefix(condition[0], "'") && strings.HasSuffix(condition[0], "'"),
		strings.HasPrefix(condition[2], "\"") && strings.HasSuffix(condition[2], "\"") || strings.HasPrefix(condition[2], "'") && strings.HasSuffix(condition[2], "'"),
	}
	var values []string
	if !isVar[0] {
		col := getColumnIndex(condition[0], cols)
		if col == -1 {
			return false, errors.New("invalid condition: column " + condition[0] + " does not exist")
		}
		values = append(values, rows[row][col].(string))
	} else {
		values = append(values, trim(condition[0]))
	}
	if !isVar[1] {
		col := getColumnIndex(condition[2], cols)
		if col == -1 {
			return false, errors.New("invalid condition: column " + condition[2] + " does not exist")
		}
		values = append(values, rows[row][col].(string))
	} else {
		values = append(values, trim(condition[2]))
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

func checkResults(results []bool, operators []string) bool {
	if len(results) == 1 {
		return results[0]
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
		var end bool
		for i := 0; i < len(result); i++ {
			if !result[i] {
				resultsFinal = append(resultsFinal, false)
				end = true
				break
			}
		}
		if !end {
			resultsFinal = append(resultsFinal, true)
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

func trim(str string) string {
	if strings.HasPrefix(str, "\"") {
		str = strings.TrimPrefix(str, "\"")
	} else {
		str = strings.TrimPrefix(str, "'")
	}
	if strings.HasSuffix(str, "\"") {
		str = strings.TrimSuffix(str, "\"")
	} else {
		str = strings.TrimSuffix(str, "'")
	}
	return str
}

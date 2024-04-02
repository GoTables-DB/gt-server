package gt_post

import (
	"errors"
	"fmt"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"strconv"
	"strings"
)

func Post(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	retTable := table.Table{}
	var retError error

	if len(query) < 1 {
		return table.Table{}, errors.New("invalid syntax")
	}
	if db == "" {
		switch strings.ToLower(query[0]) {
		case "show":
			if len(query) != 1 {
				return table.Table{}, errors.New("invalid syntax")
			}
			dbs, err := fs.GetDBs(config.Dir)
			if err != nil {
				return table.Table{}, err
			}
			retTable, retError = simpleTable("Databases", dbs)
		case "user":
			if len(query) != 0 {
				return table.Table{}, errors.New("invalid syntax")
			}
		case "backup":
			if len(query) != 0 {
				return table.Table{}, errors.New("invalid syntax")
			}
		default:
			retError = errors.New("invalid syntax")
		}
	} else if tbl == "" {
		switch strings.ToLower(query[0]) {
		case "show":
			if len(query) != 1 {
				return table.Table{}, errors.New("invalid syntax")
			}
			tables, err := fs.GetTables(db, config.Dir)
			if err != nil {
				return table.Table{}, err
			}
			retTable, retError = simpleTable("Tables", tables)
		case "create":
			if len(query) != 1 {
				return table.Table{}, errors.New("invalid syntax")
			}
			retError = fs.AddDB(db, config.Dir)
		case "move":
			if len(query) != 2 {
				return table.Table{}, errors.New("invalid syntax")
			}
			retError = fs.MoveDB(db, query[1], config.Dir)
		case "copy":
			if len(query) != 2 {
				return table.Table{}, errors.New("invalid syntax")
			}
			retError = fs.CopyDB(db, query[1], config.Dir)
		case "delete":
			if len(query) != 1 {
				return table.Table{}, errors.New("invalid syntax")
			}
			retError = fs.DeleteDB(db, config.Dir)
		default:
			retError = errors.New("invalid syntax")
		}
	} else {
		switch strings.ToLower(query[0]) {
		case "show":
			data, err := fs.GetTable(tbl, db, config.Dir)
			if err != nil {
				return table.Table{}, err
			}
			if len(query) == 1 { // Show entire table
				retTable = data
			} else if len(query) == 2 { // Show specific columns
				retTable, retError = showTable(query[1], data, []string{})
			} else if query[2] == "where" { // Show specific columns (with condition)
				retTable, retError = showTable(query[1], data, query[3:])
			} else { // Invalid syntax
				return table.Table{}, errors.New("invalid syntax")
			}
		case "create":
			if len(query) > 1 {
				retTable, retError = makeTableWithColumns(query[1:], tbl, db, config.Dir)
			} else {
				retError = fs.AddTable(tbl, db, config.Dir)
			}
		case "move":
			if len(query) != 2 {
				return table.Table{}, errors.New("invalid syntax")
			}
			retError = fs.MoveTable(tbl, query[1], db, config.Dir)
		case "copy":
			if len(query) != 2 {
				return table.Table{}, errors.New("invalid syntax")
			}
			retError = fs.CopyTable(tbl, query[1], db, config.Dir)
		case "delete":
			if len(query) != 1 {
				return table.Table{}, errors.New("invalid syntax")
			}
		case "column":
			if len(query) < 3 {
				return table.Table{}, errors.New("invalid syntax")
			}
			switch query[1] {
			case "show":
				if len(query) != 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				cols := strings.Split(query[2], ":")
				columns := []table.Column{{Name: "Name", Type: "str", Default: ""}, {Name: "Type", Type: "str", Default: ""}, {Name: "Default", Type: "str", Default: ""}}
				for i := 0; i < len(columns); i++ {
					err := data.AddColumn(columns[i])
					if err != nil {
						return table.Table{}, err
					}
				}
				var indices []int
				for i := 0; i < len(cols); i++ {
					col := getColumnIndex(cols[i], data.GetColumns())
					if col != -1 {
						indices = append(indices, col)
					}
				}
				for i := 0; i < len(indices); i++ {
					col := data.GetColumns()[indices[i]]
					row := make(map[string]any)
					rowData := []any{col.Name, col.Type, fmt.Sprint(col.Default)}
					for j := 0; j < len(columns); j++ {
						row[columns[j].Name] = rowData[j]
					}
					err := data.AddRow(row)
					if err != nil {
						return table.Table{}, err
					}
				}
				retTable = data
			case "create":
				if len(query) != 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				colSlice := strings.Split(query[2], ":")
				col := table.Column{}
				switch len(colSlice) {
				case 2:
					col.Name = colSlice[0]
					col.Type = colSlice[1]
				case 3:
					col.Name = colSlice[0]
					col.Type = colSlice[1]
					col.Default = colSlice[2]
				default:
					return table.Table{}, errors.New("invalid syntax")
				}
				err = data.AddColumn(col)
				if err != nil {
					return table.Table{}, err
				}
				err = fs.ModifyTable(data, tbl, db, config.Dir)
				retTable, retError = data, err
			case "set":
				if len(query) != 5 {
					return table.Table{}, errors.New("invalid syntax")
				}
				switch query[2] {
				case "name":
					data, err := fs.GetTable(tbl, db, config.Dir)
					if err != nil {
						return table.Table{}, err
					}
					err = data.SetColumnName(query[3], query[4])
				case "default":
					data, err := fs.GetTable(tbl, db, config.Dir)
					if err != nil {
						return table.Table{}, err
					}
					err = data.SetColumnDefault(query[3], query[4])
				default:
					retError = errors.New("invalid syntax")
				}
			case "copy":
				if len(query) != 4 {
					return table.Table{}, errors.New("invalid syntax")
				}
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				col, err := data.GetColumn(query[2])
				if err != nil {
					return table.Table{}, err
				}
				col.Name = query[3]
				err = data.AddColumn(col)
				if err != nil {
					return table.Table{}, err
				}

			case "delete":
				if len(query) != 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				tbl, err := fs.GetTable(table, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				cols := tbl.GetColumns()
				index := getColumnIndex(query[2], cols)
				if index == -1 {
					return table.Table{}, errors.New(query[2] + " is not a valid column")
				}
				cols = append(cols[:index], cols[index+1:]...)
				tbl, err = tbl.SetColumns(cols)
				if err != nil {
					return table.Table{}, err
				}
				rows := tbl.GetRows()
				for i := 0; i < len(rows); i++ {
					rows[i] = append(rows[i], rows[i][index])
				}
				tbl, err = tbl.SetRows(rows)
				if err != nil {
					return table.Table{}, err
				}
				retError = fs.ModifyTable(tbl, table, db, config.Dir)
			default:
				retError = errors.New("invalid syntax")
			}
		case "row":
			if len(query) < 3 {
				return table.Table{}, errors.New("invalid syntax")
			}
			switch query[1] {
			case "show":
				if len(query) != 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				tbl, err := fs.GetTable(table, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				rows := tbl.GetRows()
				indicesSlice := strings.Split(query[2], ":")
				indices := make([]int, 0)
				for i := 0; i < len(indicesSlice); i++ {
					index, err := strconv.Atoi(indicesSlice[i])
					if err != nil {
						return table.Table{}, err
					}
					if index < 1 || index > len(rows) {
						return table.Table{}, errors.New("index " + indicesSlice[i] + " is out of range")
					}
					indices = append(indices, index-1)
				}
				rowsNew := make([][]any, 0)
				for i := 0; i < len(indices); i++ {
					rowsNew = append(rowsNew, rows[indices[i]])
				}
				retTable, retError = tbl.SetRows(rowsNew)
			case "create":
				if len(query) < 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				tbl, err := fs.GetTable(table, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				rows := tbl.GetRows()
				rowSlice := strings.Split(query[2], ":")
				row := make([]any, 0)
				for i := 0; i < len(rowSlice); i++ {
					row = append(row, rowSlice[i])
				}
				rows = append(rows, row)
				retTable, retError = tbl.SetRows(rows)
			case "copy":
				if len(query) < 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				tbl, err := fs.GetTable(table, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				rows := tbl.GetRows()
				index, err := strconv.Atoi(query[2])
				if err != nil {
					return table.Table{}, err
				}
				if index < 1 || index > len(rows) {
					return table.Table{}, errors.New("index " + query[2] + " is out of range")
				}
				rows = append(rows, rows[index-1])
				retTable, retError = tbl.SetRows(rows)
			case "delete":
				if len(query) < 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				tbl, err := fs.GetTable(table, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				rows := tbl.GetRows()
				index, err := strconv.Atoi(query[2])
				if err != nil {
					return table.Table{}, err
				}
				if index < 1 || index > len(rows) {
					return table.Table{}, errors.New("index " + query[2] + " is out of range")
				}
				rows = append(rows[:index-1], rows[index:]...)
				retTable, retError = tbl.SetRows(rows)
			case "column": // Select cell
				retTable, retError = table.Table{}, errors.New("operations on cells not implemented yet")
			default:
				retError = errors.New("invalid syntax")
			}
		default:
			retError = errors.New("invalid syntax")
		}
	}

	return retTable, retError
}

func makeTableWithColumns(columns []string, table string, db string, dir string) (table.Table, error) {
	err := fs.AddTable(table, db, dir)
	if err != nil {
		return table.Table{}, err
	}
	tbl, err := shared.MakeTableWithColumns(columns)
	if err != nil {
		return table.Table{}, err
	}
	err = fs.ModifyTable(tbl, table, db, dir)
	return tbl, err
}

// Used to display names of databases or names of tables in a db
func simpleTable(colName string, rows []string) (table.Table, error) {
	columns := []fs.Column{{Name: colName, Type: "str"}}
	rowSlice := make([][]interface{}, 0)
	for _, row := range rows {
		rowSlice = append(rowSlice, []interface{}{row})
	}
	return shared.MakeTable(columns, rowSlice)
}

func showTable(columns string, table table.Table, condition []string) (table.Table, error) {
	rows := make([]int, 0)
	for i := 0; i < len(table.GetRows()); i++ {
		rows = append(rows, i)
	}
	if len(condition) != 0 {
		var indices []int
		if (len(condition)+1)%4 != 0 {
			return table.Table{}, errors.New("invalid condition")
		}
		for i := 0; i < len(table.GetRows()); i++ {
			var results []bool
			for j := 0; j < (len(condition)+1)/4; j++ {
				result, err := checkCondition(i, table, condition[j*4:3+j*4])
				if err != nil {
					return table.Table{}, err
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
		return table.Table{}, err
	}
	retTable, err := shared.MakeTableFromTable(cols, rows, table)
	return retTable, err
}

func checkCondition(row int, table table.Table, condition []string) (bool, error) {
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

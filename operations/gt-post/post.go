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
				tblNew := table.Table{}
				columns := []table.Column{{Name: "Name", Type: "str", Default: ""}, {Name: "Type", Type: "str", Default: ""}, {Name: "Default", Type: "str", Default: ""}}
				for i := 0; i < len(columns); i++ {
					err := tblNew.AddColumn(columns[i])
					if err != nil {
						return table.Table{}, err
					}
				}
				for i := 0; i < len(cols); i++ {
					col, err := data.GetColumn(cols[i])
					if err != nil {
						return table.Table{}, err
					}
					row := make(map[string]any)
					rowData := []any{col.Name, col.Type, fmt.Sprint(col.Default)}
					for j := 0; j < len(columns); j++ {
						row[columns[j].Name] = rowData[j]
					}
					err = tblNew.AddRow(row)
					if err != nil {
						return table.Table{}, err
					}
				}
				retTable = tblNew
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
				retError = fs.ModifyTable(data, tbl, db, config.Dir)
			case "delete":
				if len(query) != 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				err = data.DeleteColumn(query[2])
				if err != nil {
					return table.Table{}, err
				}
				retError = fs.ModifyTable(data, tbl, db, config.Dir)
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
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				rowIndices := strings.Split(query[2], ":")
				cols := data.GetColumns()
				columnIndices := make([]int, 0)
				for i := 0; i < len(cols); i++ {
					columnIndices = append(columnIndices, i)
				}
				tblNew, err := shared.MakeTableFromTable(columnIndices, []int{}, data)
				for i := 0; i < len(rowIndices); i++ {
					index, err := strconv.Atoi(rowIndices[i])
					if err != nil {
						return table.Table{}, err
					}
					row, err := data.GetRow(index)
					if err != nil {
						return table.Table{}, err
					}
					err = tblNew.AddRow(row)
					if err != nil {
						return table.Table{}, err
					}
				}
				retTable = tblNew
			case "create":
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				rowSlice := make([][]string, 0)
				for i := 2; i < len(query); i++ {
					rowSlice = append(rowSlice, strings.Split(query[i], ":"))
					if len(rowSlice[i]) != 2 {
						return table.Table{}, errors.New("invalid syntax")
					}
				}
				row := map[string]any{}
				for i := 0; i < len(rowSlice); i++ {
					row[rowSlice[i][0]] = rowSlice[i][1]
				}
				err = data.AddRow(row)
				if err != nil {
					return table.Table{}, err
				}
				retTable, retError = data, fs.ModifyTable(data, tbl, db, config.Dir)
			case "copy":
				if len(query) < 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				index, err := strconv.Atoi(query[2])
				if err != nil {
					return table.Table{}, err
				}
				row, err := data.GetRow(index)
				if err != nil {
					return table.Table{}, err
				}
				retTable, retError = data, data.AddRow(row)
			case "delete":
				if len(query) < 3 {
					return table.Table{}, errors.New("invalid syntax")
				}
				data, err := fs.GetTable(tbl, db, config.Dir)
				if err != nil {
					return table.Table{}, err
				}
				index, err := strconv.Atoi(query[2])
				if err != nil {
					return table.Table{}, err
				}
				retTable, retError = data, data.DeleteRow(index)
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

func makeTableWithColumns(columns []string, tbl string, db string, dir string) (table.Table, error) {
	err := fs.AddTable(tbl, db, dir)
	if err != nil {
		return table.Table{}, err
	}
	cols := make([]table.Column, 0)
	for i, column := range columns {
		col := table.Column{}
		switch strings.Count(column, ":") {
		case 0:
			return table.Table{}, errors.New("need to specify datatype of column at index " + strconv.Itoa(i))
		case 1:
			colSplit := strings.Split(column, ":")
			if len(colSplit) != 2 {
				return table.Table{}, errors.New("internal server error")
			}
			col.Name = colSplit[0]
			col.Type = colSplit[1]
			col.Default = ""
		case 2:
			colSplit := strings.Split(column, ":")
			if len(colSplit) != 3 {
				return table.Table{}, errors.New("internal server error")
			}
			col.Name = colSplit[0]
			col.Type = colSplit[1]
			col.Default = colSplit[2]
		default:
			return table.Table{}, errors.New("illegal column at index " + strconv.Itoa(i))
		}
		cols = append(cols, col)
	}
	tblU := table.TableU{
		Columns: cols,
		Rows:    []map[string]any{},
	}
	data, err := tblU.ToT()
	if err != nil {
		return table.Table{}, err
	}
	err = fs.ModifyTable(data, tbl, db, dir)
	return data, err
}

// Used to display names of databases or of tables in a db
func simpleTable(colName string, rowSlice []string) (table.Table, error) {
	column := table.Column{Name: colName, Type: "str"}
	data := table.Table{}
	err := data.AddColumn(column)
	if err != nil {
		return table.Table{}, err
	}
	for i := 0; i < len(rowSlice); i++ {
		row := make(map[string]any)
		row[colName] = rowSlice[i]
		err := data.AddRow(row)
		if err != nil {
			return table.Table{}, err
		}
	}
	return data, nil
}

func showTable(columns string, data table.Table, condition []string) (table.Table, error) {
	rows := make([]int, 0)
	for i := 0; i < len(data.GetRows()); i++ {
		rows = append(rows, i)
	}
	if len(condition) != 0 {
		var indices []int
		if (len(condition)+1)%4 != 0 {
			return table.Table{}, errors.New("invalid condition")
		}
		for i := 0; i < len(data.GetRows()); i++ {
			var results []bool
			for j := 0; j < (len(condition)+1)/4; j++ {
				result, err := checkCondition(i, data, condition[j*4:3+j*4])
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
	cols, err := shared.SelectColumns(colSlice, data)
	if err != nil {
		return table.Table{}, err
	}
	retTable, err := shared.MakeTableFromTable(cols, rows, data)
	return retTable, err
}

func checkCondition(rowIndex int, data table.Table, condition []string) (bool, error) {
	if len(condition) != 3 {
		return false, errors.New("invalid condition")
	}
	row, err := data.GetRow(rowIndex)
	if err != nil {
		return false, err
	}
	isVar := []bool{
		strings.HasPrefix(condition[0], "\"") && strings.HasSuffix(condition[0], "\"") || strings.HasPrefix(condition[0], "'") && strings.HasSuffix(condition[0], "'"),
		strings.HasPrefix(condition[2], "\"") && strings.HasSuffix(condition[2], "\"") || strings.HasPrefix(condition[2], "'") && strings.HasSuffix(condition[2], "'"),
	}
	var values []string
	if !isVar[0] {
		col, err := data.GetColumn(condition[0])
		if err != nil {
			return false, errors.New("invalid condition: column " + condition[0] + " does not exist")
		}
		values = append(values, row[col.Name].(string))
	} else {
		values = append(values, trim(condition[0]))
	}
	if !isVar[1] {
		col, err := data.GetColumn(condition[2])
		if err != nil {
			return false, errors.New("invalid condition: column " + condition[2] + " does not exist")
		}
		values = append(values, row[col.Name].(string))
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

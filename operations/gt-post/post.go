package gt_post

import (
	"errors"
	"fmt"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"strconv"
	"strings"
	"time"
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
			retTable, retError = rootShow(query, config)
		case "user":
			retError = errors.New("users not implemented yet")
		case "backup":
			retError = errors.New("backups not implemented yet")
		default:
			retError = errors.New("invalid syntax")
		}
	} else if tbl == "" {
		switch strings.ToLower(query[0]) {
		case "show":
			retTable, retError = dbShow(query, db, config)
		case "create":
			retTable, retError = dbCreate(query, db, config)
		case "set":
			if len(query) < 3 {
				return table.Table{}, errors.New("invalid syntax")
			}
			switch strings.ToLower(query[1]) {
			case "name":
				retTable, retError = dbSetName(query, db, config)
			default:
				retError = errors.New("invalid syntax")
			}
		case "copy":
			retTable, retError = dbCopy(query, db, config)
		case "delete":
			retTable, retError = dbDelete(query, db, config)
		default:
			retError = errors.New("invalid syntax")
		}
	} else {
		switch strings.ToLower(query[0]) {
		case "show":
			retTable, retError = tableShow(query, tbl, db, config)
		case "create":
			retTable, retError = tableCreate(query, tbl, db, config)
		case "set":
			if len(query) < 3 {
				return table.Table{}, errors.New("invalid syntax")
			}
			switch strings.ToLower(query[1]) {
			case "name":
				retTable, retError = tableSetName(query, tbl, db, config)
			default:
				retError = errors.New("invalid syntax")
			}
		case "copy":
			retTable, retError = tableCopy(query, tbl, db, config)
		case "delete":
			retTable, retError = tableDelete(query, tbl, db, config)
		case "column":
			if len(query) < 3 {
				return table.Table{}, errors.New("invalid syntax")
			}
			switch query[1] {
			case "show":
				retTable, retError = columnShow(query, tbl, db, config)
			case "create":
				retTable, retError = columnCreate(query, tbl, db, config)
			case "set":
				switch query[2] {
				case "name":
					retTable, retError = columnSetName(query, tbl, db, config)
				case "default":
					retTable, retError = columnSetDefault(query, tbl, db, config)
				default:
					retError = errors.New("invalid syntax")
				}
			case "copy":
				retTable, retError = columnCopy(query, tbl, db, config)
			case "delete":
				retTable, retError = columnDelete(query, tbl, db, config)
			default:
				retError = errors.New("invalid syntax")
			}
		case "row":
			if len(query) < 2 {
				return table.Table{}, errors.New("invalid syntax")
			}
			switch query[1] {
			case "show":
				retTable, retError = rowShow(query, tbl, db, config)
			case "create":
				retTable, retError = rowCreate(query, tbl, db, config)
			case "set": // Select cell
				retTable, retError = rowSet(query, tbl, db, config)
			case "copy":
				retTable, retError = rowCopy(query, tbl, db, config)
			case "delete":
				retTable, retError = rowDelete(query, tbl, db, config)
			default:
				retError = errors.New("invalid syntax")
			}
		default:
			retError = errors.New("invalid syntax")
		}
	}
	return retTable, retError
}

func rootShow(query []string, config fs.Conf) (table.Table, error) {
	if len(query) != 1 {
		return table.Table{}, errors.New("invalid syntax")
	}
	dbs, err := fs.GetDBs(config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	return simpleTable("Databases", dbs)
}

func dbShow(query []string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 1 {
		return table.Table{}, errors.New("invalid syntax")
	}
	tables, err := fs.GetTables(db, config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	return simpleTable("Tables", tables)
}

func dbCreate(query []string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 1 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.AddDB(db, config.Dir)
}

func dbSetName(query []string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 3 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.RenameDB(db, query[2], config.Dir)
}

func dbCopy(query []string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 2 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.CopyDB(db, query[1], config.Dir)
}

func dbDelete(query []string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 1 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.DeleteDB(db, config.Dir)
}

func tableShow(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	data, err := fs.GetTable(tbl, db, config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	if len(query) == 1 { // Show entire table
		return data, nil
	} else if len(query) == 2 { // Show specific columns
		return showTable(query[1], data, []string{})
	} else if query[2] == "where" { // Show specific columns (with condition)
		return showTable(query[1], data, query[3:])
	} else { // Invalid syntax
		return table.Table{}, errors.New("invalid syntax")
	}
}

func tableCreate(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) > 1 {
		return makeTableWithColumns(query[1:], tbl, db, config.Dir)
	} else {
		return table.Table{}, fs.AddTable(tbl, db, config.Dir)
	}
}

func tableSetName(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 3 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.RenameTable(tbl, query[2], db, config.Dir)
}

func tableCopy(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 2 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.CopyTable(tbl, query[1], db, config.Dir)
}

func tableDelete(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 1 {
		return table.Table{}, errors.New("invalid syntax")
	}
	return table.Table{}, fs.DeleteTable(tbl, db, config.Dir)
}

func columnShow(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
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
	return tblNew, nil
}

func columnCreate(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
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
	return data, err
}

func columnSetName(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 5 {
		return table.Table{}, errors.New("invalid syntax")
	}
	data, err := fs.GetTable(tbl, db, config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	err = data.SetColumnName(query[3], query[4])
	if err != nil {
		return table.Table{}, err
	}
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func columnSetDefault(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 5 {
		return table.Table{}, errors.New("invalid syntax")
	}
	data, err := fs.GetTable(tbl, db, config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	err = data.SetColumnDefault(query[3], query[4])
	if err != nil {
		return table.Table{}, err
	}
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func columnCopy(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
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
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func columnDelete(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
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
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func rowShow(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
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
	if err != nil {
		return table.Table{}, err
	}
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
	return tblNew, nil
}

func rowCreate(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) < 2 {
		return table.Table{}, errors.New("invalid syntax")
	}
	data, err := fs.GetTable(tbl, db, config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	rowSlice := make([][]string, 0)
	for i := 2; i < len(query); i++ {
		rowSlice = append(rowSlice, strings.Split(query[i], ":"))
		if len(rowSlice[i-2]) != 2 {
			return table.Table{}, errors.New("invalid syntax")
		}
	}
	row := map[string]any{}
	for i := 0; i < len(rowSlice); i++ {
		col, err := data.GetColumn(rowSlice[i][0])
		if err != nil {
			return table.Table{}, err
		}
		cell, err := convert(rowSlice[i][1], col.Type)
		if err != nil {
			return table.Table{}, err
		}
		row[rowSlice[i][0]] = cell
	}
	err = data.AddRow(row)
	if err != nil {
		return table.Table{}, err
	}
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func rowSet(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 4 {
		return table.Table{}, errors.New("invalid syntax")
	}
	data, err := fs.GetTable(tbl, db, config.Dir)
	if err != nil {
		return table.Table{}, err
	}
	indices := strings.Split(query[2], ":")
	if len(indices) != 2 {
		return table.Table{}, errors.New("invalid syntax")
	}
	index, err := strconv.Atoi(indices[0])
	if err != nil {
		return table.Table{}, err
	}
	row, err := data.GetRow(index)
	if err != nil {
		return table.Table{}, err
	}
	col, err := data.GetColumn(indices[1])
	if err != nil {
		return table.Table{}, err
	}
	cell, err := convert(query[3], col.Type)
	if err != nil {
		return table.Table{}, err
	}
	row[indices[1]] = cell
	err = data.SetRow(index, row)
	if err != nil {
		return table.Table{}, err
	}
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func rowCopy(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 3 {
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
	err = data.AddRow(row)
	if err != nil {
		return table.Table{}, err
	}
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
}

func rowDelete(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	if len(query) != 3 {
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
	err = data.DeleteRow(index)
	if err != nil {
		return table.Table{}, err
	}
	return data, fs.ModifyTable(data, tbl, db, config.Dir)
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

func convert(value string, datatype string) (any, error) {
	switch datatype {
	// String
	case "str":
		return value, nil
	// Number
	case "num", "int", "flt":
		ret, err := strconv.ParseFloat(value, 64)
		return ret, err
	// Boolean
	case "bol":
		ret, err := strconv.ParseBool(value)
		return ret, err
	// Date
	case "dat":
		ret, err := time.Parse(time.RFC3339, value)
		return ret, err
	// Table
	case "tbl":
		return nil, errors.New("not implemented yet")
	default:
		return nil, errors.New("invalid datatype")
	}
}

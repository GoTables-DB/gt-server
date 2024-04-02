package shared

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"log"
)

/// Table operations ///

func MakeTable(columns []table.Column, rows []map[string]any) (table.Table, error) {
	tblU := table.TableU{}
	tblU.Columns = columns
	tblU.Rows = rows
	tbl, err := tblU.ToT()
	return tbl, err
}

/*
func MakeTableWithColumns(columns []string) (table.Table, error) {
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
	tbl, err := tblU.ToT()
	return tbl, err
}
*/

func MakeTableFromTable(columnIndices []int, rowIndices []int, tbl table.Table) (table.Table, error) {
	if len(columnIndices) == 0 {
		return table.Table{}, nil
	}
	tblU := table.TableU{}
	columns := make([]table.Column, len(columnIndices))
	columnsOld := tbl.GetColumns()
	for i, j := range columnIndices {
		columns[i].Name = columnsOld[j].Name
		columns[i].Type = columnsOld[j].Type
	}
	tblU.Columns = columns
	rows := make([]map[string]any, len(rowIndices))
	rowsOld := tbl.GetRows()
	for i, j := range rowIndices {
		row := make(map[string]any, len(columnIndices))
		for _, k := range columnIndices {
			row[columnsOld[k].Name] = rowsOld[j][columnsOld[k].Name]
		}
		rows[i] = row
	}
	tblU.Rows = rows
	tbl, err := tblU.ToT()
	return tbl, err
}

func SelectColumns(columnNames []string, tbl table.Table) ([]int, error) {
	if len(columnNames) == 0 {
		return nil, errors.New("no columns specified")
	}
	indices := make([]int, 0)
	for _, columnName := range columnNames {
		if columnName == "*" {
			indicesAll := make([]int, 0)
			for i := 0; i < len(tbl.GetColumns()); i++ {
				indicesAll = append(indicesAll, i)
			}
			return indicesAll, nil
		}
		var inTable bool
		for i, column := range tbl.GetColumns() {
			if columnName == column.Name {
				indices = append(indices, i)
				inTable = true
				break
			}
		}
		if !inTable {
			return []int{}, errors.New("column " + columnName + " does not exist")
		}
	}
	return indices, nil
}

/// Filesystem operations ///

func AddDB(db, dir string) error {
	if db == "" {
		return errors.New("no database specified")
	}
	dbs, dbErr := fs.GetDBs(dir)
	if dbErr != nil {
		log.Println(dbErr)
		return dbErr
	}
	for _, name := range dbs {
		if db == name {
			return errors.New("database already exists")
		}
	}
	return fs.AddDB(db, dir)
}

func AddTable(table, db, dir string) error {
	if db == "" {
		return errors.New("no database specified")
	}
	if table == "" {
		return errors.New("no table specified")
	}
	tables, tableErr := fs.GetTables(db, dir)
	if tableErr != nil {
		log.Println(tableErr)
		return tableErr
	}
	for _, name := range tables {
		if table == name {
			return errors.New("table already exists")
		}
	}
	return fs.AddTable(table, db, dir)
}

func DeleteDB(db, dir string) error {
	if db == "" {
		return errors.New("no database specified")
	}
	dbs, dbErr := fs.GetDBs(dir)
	if dbErr != nil {
		log.Println(dbErr)
		return dbErr
	}
	for _, name := range dbs {
		if db == name {
			return fs.DeleteDB(db, dir)
		}
	}
	return errors.New("database not found")
}

func DeleteTable(table, db, dir string) error {
	if db == "" {
		return errors.New("no database specified")
	}
	if table == "" {
		return errors.New("no table specified")
	}
	tables, tableErr := fs.GetTables(db, dir)
	if tableErr != nil {
		log.Println(tableErr)
		return tableErr
	}
	for _, name := range tables {
		if table == name {
			return fs.DeleteTable(table, db, dir)
		}
	}
	return errors.New("table not found")
}

/// User operations ///

func AddUser() {

}

func ModifyUser() {

}

func DeleteUser() {

}

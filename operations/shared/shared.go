package shared

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"log"
	"strconv"
	"strings"
)

func MakeTable(columns []fs.Column, rows [][]interface{}) (fs.Table, error) {
	table := fs.Table{}
	var err error = nil
	table = table.SetColumns(columns)
	table, err = table.SetRows(rows)
	return table, err
}

func MakeTableWithColumns(columns []string) (fs.Table, error) {
	cols := make([]fs.Column, 0)
	for i, column := range columns {
		col := fs.Column{}
		switch strings.Count(column, ":") {
		case 0:
			return fs.Table{}, errors.New("need to specify datatype of column at index " + strconv.Itoa(i))
		case 1:
			colSplit := strings.Split(column, ":")
			if len(colSplit) != 2 {
				return fs.Table{}, errors.New("internal server error")
			}
			datatype, err := fs.DetermineDatatype(colSplit[1])
			if err != nil {
				return fs.Table{}, err
			}
			col.Name = colSplit[0]
			col.Type = datatype
			col.Default = datatype
		case 2:
			colSplit := strings.Split(column, ":")
			if len(colSplit) != 3 {
				return fs.Table{}, errors.New("internal server error")
			}
			datatype, err := fs.DetermineDatatype(colSplit[1])
			if err != nil {
				return fs.Table{}, err
			}
			col.Name = colSplit[0]
			col.Type = datatype
			col.Default = colSplit[2]
		default:
			return fs.Table{}, errors.New("illegal column at index " + strconv.Itoa(i))
		}
		cols = append(cols, col)
	}
	tbl := fs.Table{}.SetColumns(cols)
	return tbl, nil
}

func MakeTableFromTable(columnIndices []int, rowIndices []int, table fs.Table) (fs.Table, error) {
	retTable := fs.Table{}
	if len(columnIndices) == 0 {
		return fs.Table{}, nil
	}
	retTable = retTable.SetColumns(make([]fs.Column, len(columnIndices)))
	columns := retTable.GetColumns()
	columnsOld := table.GetColumns()
	for i, j := range columnIndices {
		columns[i].Name = columnsOld[j].Name
		columns[i].Type = columnsOld[j].Type
	}
	retTable = retTable.SetColumns(columns)
	rows := make([][]interface{}, len(rowIndices))
	rowsOld := table.GetRows()
	for _, i := range rowIndices {
		for _, j := range columnIndices {
			rows[i] = append(rows[i], rowsOld[i][j])
		}
	}
	return retTable.SetRows(rows)
}

func SelectTable(tableName string, db string, config fs.Conf) (fs.Table, error) {
	if tableName == "" {
		return fs.Table{}, errors.New("no table specified")
	}
	table, tableErr := fs.GetTable(db, tableName, config.Dir)
	if tableErr != nil {
		log.Println(tableErr)
		return fs.Table{}, tableErr
	}
	return table, errors.New("table not found")
}

func SelectColumns(columnNames []string, table fs.Table) ([]int, error) {
	if len(columnNames) == 0 {
		return nil, errors.New("no columns specified")
	}
	Indices := make([]int, 0)
	for _, columnName := range columnNames {
		if columnName == "*" {
			for i := 0; i < len(table.GetColumns()); i++ {
				Indices = append(Indices, i)
			}
			return Indices, nil
		}
		for i, column := range table.GetColumns() {
			if columnName == column.Name {
				Indices = append(Indices, i)
				break
			}
		}
	}
	return Indices, nil
}

func ModifyRow() {

}

func ModifyColumn() {

}

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
	return fs.NewDB(db, dir)
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
	return fs.NewTable(table, db, dir)
}

func AddRow() {

}

func AddColumn() {

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

func DeleteRow() {

}

func DeleteColumn() {

}

func AddUser() {

}

func ModifyUser() {

}

func DeleteUser() {

}

package shared

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"log"
)

func MakeNewTable(columns []fs.Column, rows [][]interface{}) (fs.Table, error) {
	table := fs.Table{}
	var err error = nil
	table = table.SetColumns(columns)
	table, err = table.SetRows(rows)
	return table, err
}

func MakeTableFromTable(columnIndexes []int, rowIndexes []int, table fs.Table) (fs.Table, error) {
	retTable := fs.Table{}
	if len(columnIndexes) == 0 {
		return fs.Table{}, nil
	}
	retTable = retTable.SetColumns(make([]fs.Column, len(table.GetColumns())))
	for i, column := range table.GetColumns() {
		if column.Name == "*" {
			retTable = retTable.SetColumns(table.GetColumns())
			break
		}
		columns := retTable.GetColumns()
		columns[i].Name = column.Name
		columns[i].Type = column.Type
		retTable = retTable.SetColumns(columns)
	}
	rows := table.GetRows()
	rowsNew := make([][]interface{}, len(rowIndexes))
	for i := range rowIndexes {
		for j := range columnIndexes {
			rowsNew[i] = append(rowsNew[i], rows[i][j])
		}
	}
	return retTable.SetRows(rowsNew)
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

func SelectRows() {

}

func SelectColumns(columnNames []string, table fs.Table) ([]int, error) {
	if len(columnNames) == 0 {
		return nil, errors.New("no columns specified")
	}
	indexes := make([]int, 0)
	for _, columnName := range columnNames {
		if columnName == "*" {
			for i := 0; i < len(table.GetColumns()); i++ {
				indexes = append(indexes, i)
			}
			return indexes, nil
		}
		for i, column := range table.GetColumns() {
			if columnName == column.Name {
				indexes = append(indexes, i)
			}
		}
	}
	return indexes, nil
}

func ModifyDB() {

}

func ModifyTable() {

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

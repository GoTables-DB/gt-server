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
	retTable, err := retTable.SetRows(rowsNew)
	return retTable, err
}

func SelectTable(db string, tableName string, config fs.Conf) (fs.Table, error) {
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

func AddDB() {

}

func AddTable() {

}

func AddRow() {

}

func AddColumn() {

}

func DeleteDB() {

}

func DeleteTable() {

}

func DeleteRow() {

}

func DeleteColumn() {

}

func DddUser() {

}

func ModifyUser() {

}

func DeleteUser() {

}

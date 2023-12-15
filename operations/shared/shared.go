package shared

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"log"
)

func MakeNewTable(columns []fs.Column, rows [][]interface{}) fs.Table {
	return fs.Table{
		ColumnNames: columns,
		Rows:        rows,
	}
}

func MakeTableFromExisting(columnIndexes []int, rowIndexes []int, table fs.Table) fs.Table {
	retTable := fs.Table{}
	if len(columnIndexes) == 0 {
		return fs.Table{}
	}
	retTable.ColumnNames = make([]fs.Column, len(table.ColumnNames))
	for i, column := range table.ColumnNames {
		if column.Name == "*" {
			retTable.ColumnNames = table.ColumnNames
			break
		}
		retTable.ColumnNames[i].Name = column.Name
		retTable.ColumnNames[i].Type = column.Type
	}
	if len(rowIndexes) == 0 {
		return fs.Table{
			ColumnNames: retTable.ColumnNames,
			Rows:        nil,
		}
	}
	for i := range rowIndexes {
		retTable.Rows = append(retTable.Rows, make([]interface{}, len(columnIndexes)))
		for j := range columnIndexes {
			retTable.Rows[i] = append(retTable.Rows[i], table.Rows[i][j])
		}
	}
	return retTable
}

func selectTable(db string, tableName string, config fs.Conf) (fs.Table, error) {
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

func selectRows() {

}

func selectColumns(columnNames []string, table fs.Table) ([]int, error) {
	if len(columnNames) == 0 {
		return nil, errors.New("no columns specified")
	}
	indexes := make([]int, 0)
	for _, columnName := range columnNames {
		if columnName == "*" {
			for i := 0; i < len(table.ColumnNames); i++ {
				indexes = append(indexes, i)
			}
			return indexes, nil
		}
		for i, column := range table.ColumnNames {
			if columnName == column.Name {
				indexes = append(indexes, i)
			}
		}
	}
	return indexes, nil
}

func modifyDB() {

}

func modifyTable() {

}

func modifyRow() {

}

func modifyColumn() {

}

func AddDB() {

}

func AddTable() {

}

func addRow() {

}

func addColumn() {

}

func deleteDB() {

}

func deleteTable() {

}

func deleteRow() {

}

func deleteColumn() {

}

func addUser() {

}

func modifyUser() {

}

func deleteUser() {

}

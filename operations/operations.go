package operations

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-delete"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-get"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-post"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-put"
	"git.jereileu.ch/gotables/server/gt-server/operations/sql-post"
	"log"
	"net/http"
	"strings"
)

/// Incoming requests ///

// Handle requests that use GoTables syntax

func GTSyntax(method string, dir string, query string, config fs.Conf) (fs.Table, error) {
	if !config.EnableGTSyntax {
		return fs.Table{}, errors.New("gotables syntax is disabled on the server")
	}
	query = strings.TrimSpace(query)
	querySlice := strings.Split(query, " ")
	db, table, err := dirSplit(dir)
	if err != nil {
		return fs.Table{}, err
	}
	retTable, retError := gtQuery(method, querySlice, db, table, config)
	return retTable, retError
}

// Handle requests that use SQL syntax

func SQLSyntax(method string, dir string, query string, config fs.Conf) (fs.Table, error) {
	if !config.EnableGTSyntax {
		return fs.Table{}, errors.New("sql syntax is disabled on the server")
	}
	query = strings.TrimSpace(query)
	querySlice := strings.Split(query, " ")
	db, table, err := dirSplit(dir)
	if err != nil {
		return fs.Table{}, err
	}
	retTable, retError := sqlQuery(method, querySlice, db, table, config)
	return retTable, retError
}

// / Request operations ///
func dirSplit(dir string) (string, string, error) {
	dir = strings.TrimPrefix(dir, "/")
	dir = strings.TrimSuffix(dir, "/")
	dirSlice := strings.Split(dir, "/")
	db := ""
	table := ""
	if len(dirSlice) == 0 {
		return "", "", errors.New("no database specified")
	} else if len(dirSlice) > 2 {
		return "", "", errors.New("path too long")
	} else if len(dirSlice) == 1 {
		db = dirSlice[0]
	} else {
		db = dirSlice[0]
		table = dirSlice[1]
	}
	if db == "" {
		return "", "", errors.New("no database specified")
	}
	return db, table, nil
}

func gtQuery(method string, query []string, db string, table string, config fs.Conf) (fs.Table, error) {
	if len(query) == 0 {
		return fs.Table{}, errors.New("empty query")
	}
	retTable := fs.Table{}
	var retError error = nil
	if method == http.MethodGet || method == http.MethodHead {
		retTable, retError = gt_get.Get(db, table, config)
	} else if method == http.MethodPut {
		retTable, retError = gt_put.Put(db, table, config)
	} else if method == http.MethodPost {
		retTable, retError = gt_post.Post(query, db, table, config)
	} else if method == http.MethodDelete {
		retTable, retError = gt_del.Del(db, table, config)
	} else {
		return fs.Table{}, errors.New("invalid method")
	}
	return retTable, retError
}

func sqlQuery(method string, query []string, db string, table string, config fs.Conf) (fs.Table, error) {
	if len(query) == 0 {
		return fs.Table{}, errors.New("empty query")
	}
	retTable := fs.Table{}
	var retError error = nil
	if method == http.MethodPost {
		retTable, retError = sql_post.Post(query, db, table, config)
	} else {
		return fs.Table{}, errors.New("invalid method")
	}
	return retTable, retError
}

// / Operations on db/tables ///
func login() {

}

func logout() {

}

func checkDB(dbName string, config fs.Conf) (bool, error) {
	if dbName == "" {
		return false, errors.New("no database specified")
	}
	dbs, dbErr := fs.GetDBs(config.Dir)
	if dbErr != nil {
		log.Println(dbErr)
		return false, dbErr
	}
	for _, db := range dbs {
		if dbName == db {
			return true, nil
		}
	}
	return false, nil
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

func MakeTableNew(columns []fs.Column, rows [][]interface{}) fs.Table {
	return fs.Table{
		ColumnNames: columns,
		Rows:        rows,
	}
}

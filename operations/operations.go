package operations

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"log"
)

/// Incoming requests ///

// Handle requests that use GoTables syntax

func GTSyntax(table fs.Table, operation string, config fs.Conf) (fs.Table, error) {
	if config.EnableGTSyntax {
		return fs.Table{}, nil
	}
	return fs.Table{}, errors.New("gotables syntax is disabled on the server")
}

// Handle requests that use SQL syntax

func SQLsyntax(table fs.Table, operation string, config fs.Conf) (fs.Table, error) {
	if config.EnableSQLSyntax {
		return fs.Table{}, nil
	}
	return fs.Table{}, errors.New("sql syntax is disabled on the server")
}

// Handle requests that use go syntax to directly use functions

func GoSyntax(table fs.Table, operation string, config fs.Conf) (fs.Table, error) {
	if config.EnableGoSyntax {
		return fs.Table{}, nil
	}
	return fs.Table{}, errors.New("go syntax is disabled on the server")
}

// / Operations on tables ///
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

func selectColumns() {

}

func modifyDB() {

}

func modifyTable() {

}

func modifyRow() {

}

func modifyColumn() {

}

func addDB() {

}

func addTable() {

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

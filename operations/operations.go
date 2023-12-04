package operations

import "git.jereileu.ch/gotables/server/gt-server/fs"

/// Incoming requests ///

// Handle requests that use GoTables syntax

func GTSyntax(table fs.Table, operation string) fs.Table {
	return fs.Table{}
}

// Handle requests that use SQL syntax

func SQLsyntax(table fs.Table, operation string) fs.Table {
	return fs.Table{}
}

// Handle requests that use go syntax to directly use functions

func ExecFunc(table fs.Table, operation string) fs.Table {
	return fs.Table{}
}

// / Operations on tables ///
func login() {

}

func logout() {

}

func selectDB() {

}

func selectTable() {

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

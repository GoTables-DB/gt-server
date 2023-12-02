package operations

import "github.com/Jero075/gotables/fs"

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

/// Operations on tables ///

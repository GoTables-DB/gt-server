package gt_post

import (
	"errors"
	"fmt"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-get"
	"strings"
)

func Post(query []string, table string, db string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error

	fmt.Println(query)

	switch strings.ToLower(query[0]) {
	// List dbs or tables
	case "list":
		if len(query) != 2 {
			retError = errors.New("invalid syntax")
		} else {
			switch strings.ToLower(query[1]) {
			case "databases":
				if db != "" {
					retError = errors.New("invalid syntax")
				} else {
					retTable, retError = gt_get.Get("", "", config)
				}
			case "tables":
				if table == "" || db != "" {
					retTable, retError = gt_get.Get("", db, config)
				} else {
					retError = errors.New("invalid syntax")
				}
			}
		}
	// Create db or table
	case "create":
	// Insert column or row
	case "insert":
	// Modify db or table
	case "modify":
	// Change row content
	case "change":
	// Delete db, table, column or row
	case "delete":
	}

	return retTable, retError
}

package gt_put

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
)

func Put(db string, table string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error = nil

	if db == "" {
		retError = errors.New("no database specified")
	} else if table == "" {
		shared.AddDB()
	} else {
		shared.AddTable()
	}

	return retTable, retError
}

package gt_put

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
)

func Put(table string, db string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error = nil

	if db == "" {
		retError = errors.New("no database specified")
	} else if table == "" {
		retError = shared.AddDB(db, config.Dir)
	} else {
		retError = shared.AddTable(table, db, config.Dir)
	}

	return retTable, retError
}

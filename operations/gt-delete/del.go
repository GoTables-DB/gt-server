package gt_del

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/shared"
	"git.jereileu.ch/gotables/server/gt-server/table"
)

func Del(tbl string, db string, config fs.Conf) (table.Table, error) {
	retTable := table.Table{}
	var retError error = nil

	if db == "" {
		retError = errors.New("no database specified")
	} else if tbl == "" {
		retError = shared.DeleteDB(db, config.Dir)
	} else {
		retError = shared.DeleteTable(tbl, db, config.Dir)
	}

	return retTable, retError
}

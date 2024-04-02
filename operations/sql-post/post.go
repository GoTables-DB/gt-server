package sql_post

import (
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/table"
)

func Post(query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	retTable := table.Table{}
	var retError error = nil
	return retTable, retError
}

package sql_post

import "git.jereileu.ch/gotables/server/gt-server/fs"

func Post(query []string, table string, db string, config fs.Conf) (fs.Table, error) {
	retTable := fs.Table{}
	var retError error = nil
	return retTable, retError
}

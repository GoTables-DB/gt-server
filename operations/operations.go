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
	table, db, err := dirSplit(dir)
	if err != nil {
		return fs.Table{}, err
	}
	retTable, retError := gtQuery(method, querySlice, table, db, config)
	return retTable, retError
}

// Handle requests that use SQL syntax

func SQLSyntax(method string, dir string, query string, config fs.Conf) (fs.Table, error) {
	if !config.EnableGTSyntax {
		return fs.Table{}, errors.New("sql syntax is disabled on the server")
	}
	query = strings.TrimSpace(query)
	querySlice := strings.Split(query, " ")
	table, db, err := dirSplit(dir)
	if err != nil {
		return fs.Table{}, err
	}
	retTable, retError := sqlQuery(method, querySlice, table, db, config)
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
	return table, db, nil
}

func gtQuery(method string, query []string, table string, db string, config fs.Conf) (fs.Table, error) {
	if len(query) == 0 {
		return fs.Table{}, errors.New("empty query")
	}
	retTable := fs.Table{}
	var retError error = nil
	if method == http.MethodGet || method == http.MethodHead {
		retTable, retError = gt_get.Get(table, db, config)
	} else if method == http.MethodPut {
		retTable, retError = gt_put.Put(table, db, config)
	} else if method == http.MethodPost {
		retTable, retError = gt_post.Post(query, table, db, config)
	} else if method == http.MethodDelete {
		retTable, retError = gt_del.Del(table, db, config)
	} else {
		return fs.Table{}, errors.New("invalid method")
	}
	return retTable, retError
}

func sqlQuery(method string, query []string, table string, db string, config fs.Conf) (fs.Table, error) {
	if len(query) == 0 {
		return fs.Table{}, errors.New("empty query")
	}
	retTable := fs.Table{}
	var retError error = nil
	if method == http.MethodPost {
		retTable, retError = sql_post.Post(query, table, db, config)
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

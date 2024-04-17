package operations

import (
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-delete"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-get"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-post"
	"git.jereileu.ch/gotables/server/gt-server/operations/gt-put"
	"git.jereileu.ch/gotables/server/gt-server/operations/sql-post"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"html"
	"net/http"
	"strings"
)

func GTSyntax(method string, dir string, query string, config fs.Conf) (table.Table, error) {
	if !config.EnableGTSyntax {
		return table.Table{}, errors.New("gotables syntax is disabled on the server")
	}
	querySlice := getQuerySlice(query)
	tbl, db, err := dirSplit(dir)
	if err != nil {
		return table.Table{}, err
	}
	retTable, retError := gtQuery(method, querySlice, tbl, db, config)
	return retTable, retError
}

func SQLSyntax(method string, dir string, query string, config fs.Conf) (table.Table, error) {
	if !config.EnableSQLSyntax {
		return table.Table{}, errors.New("sql syntax is disabled on the server")
	}
	querySlice := getQuerySlice(query)
	tbl, db, err := dirSplit(dir)
	if err != nil {
		return table.Table{}, err
	}
	retTable, retError := sqlQuery(method, querySlice, tbl, db, config)
	return retTable, retError
}

func dirSplit(dir string) (string, string, error) {
	dir = strings.TrimPrefix(dir, "/")
	dir = strings.TrimSuffix(dir, "/")
	dirSlice := strings.Split(dir, "/")
	db := ""
	tbl := ""
	if len(dirSlice) > 2 {
		return "", "", errors.New("path too long")
	} else if len(dirSlice) == 1 {
		db = dirSlice[0]
	} else {
		db = dirSlice[0]
		tbl = dirSlice[1]
	}
	return tbl, db, nil
}

func getQuerySlice(query string) []string {
	query = html.UnescapeString(query)
	query = strings.TrimSpace(query)
	return strings.Split(query, " ")
}

func gtQuery(method string, query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	retTable := table.Table{}
	var retError error = nil
	if method == http.MethodGet || method == http.MethodHead {
		retTable, retError = gt_get.Get(tbl, db, config)
	} else if method == http.MethodPut {
		retTable, retError = gt_put.Put(tbl, db, config)
	} else if method == http.MethodPost {
		retTable, retError = gt_post.Post(query, tbl, db, config)
	} else if method == http.MethodDelete {
		retTable, retError = gt_del.Del(tbl, db, config)
	} else {
		return table.Table{}, errors.New("invalid method")
	}
	return retTable, retError
}

func sqlQuery(method string, query []string, tbl string, db string, config fs.Conf) (table.Table, error) {
	retTable := table.Table{}
	var retError error = nil
	if method == http.MethodPost {
		retTable, retError = sql_post.Post(query, tbl, db, config)
	} else {
		return table.Table{}, errors.New("invalid method")
	}
	return retTable, retError
}

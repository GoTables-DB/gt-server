package data

import "crypto/sha512"

type DB struct {
	name    string
	tables  []Table
	dbusers map[string][64]byte
}

type Table struct {
	name    string
	access  map[string]int
	indexes []map[string]int
	table   []map[string]any
}

func ls(dir string) (error error, contents []string) {
	// Print content of a directory
	return nil, nil
}

func InitDB(name string, adminUsername string, adminPassword string) error {
	pwHash := sha512.Sum512([]byte(adminPassword))
	db := DB{name: name, tables: make([]Table, 0), dbusers: make(map[string][64]byte)}
	db.dbusers[adminUsername] = pwHash

	return nil
}

func GetDBs() (error error, dbList []string) {
	err, dbs := ls(config.rootDir)
	return err, dbs
}

func GetTables(db string) (error error, tblList []string) {
	err, tbls := ls(config.rootDir + "/" + db)
	return err, tbls
}

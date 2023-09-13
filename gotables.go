package gotables

import (
	"crypto/sha256"
	"fmt"
)

type DB struct {
	name    string
	port    int
	tables  []Table
	dbusers map[string][32]byte
}

type Table struct {
	name    string
	access  map[string]int
	indexes []map[string]int
	table   []map[string]any
}

func InitDB(name string, adminUsername string, adminPassword string, port int) {
	if port == 0 {
		port = 5678
	}
	pwHash := sha256.Sum256([]byte(adminPassword))
	db := DB{name: name, port: port}
	db.dbusers["adminUsername"] = pwHash
	fmt.Print(db)
}

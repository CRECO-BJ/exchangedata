package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DataStore ...
type DataStore struct {
	db *sql.DB
}

// Start -- DataStore
func (d *DataStore) Start() error {
	db, err := sql.Open("mysql", "mysql:creco!73@/exdata")
	log.Println("DB open with", err, db.Stats())
	d.db = db
	return err
}

// Stop -- DataStore
func (d *DataStore) Stop() {
	d.db.Close()
}

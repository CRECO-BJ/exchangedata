package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type dbTrade struct {
	gorm.Model
	Trade
}

// DataStore ...
type DataStore struct {
	db *sql.DB
}

// Start -- DataStore
func (d *DataStore) Open() error {
	db, err := gorm.Open("mysql", "root:office98@/exdata")
	if err != nil {
		log.Println("DB open with", err, db.Stats())
		return err
	}
	d.db = db

	return err
}

// Stop -- DataStore
func (d *DataStore) Close() {
	d.db.Close()
}

// StoreTrade stores a trader to the database
func (d *DataStore) StoreTrade(t *Trade) error {
	dbTrade := &dbTrade{Trade: t}
	d.db.Create(dbTrade)
}

// CreateDB creates the database if it is not exist
func CreateDB() {
	name := "exdata"
	passwd := "office98"
	user := "root"

	db, err := sql.Open("mysql", user+":"+passwd+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP TABLE example")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE MYSQL")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DROP DATABASE " + name)
	if err != nil {
		panic(err)
	}
}

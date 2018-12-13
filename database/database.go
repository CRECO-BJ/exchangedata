package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DataStore ...
type DataStore struct {
	db *gorm.DB
}

func (d *DataStore) Open() error {
	db, err := gorm.Open("mysql", "root:office98@/exdata")
	if err != nil {
		log.Println("DB open with", err)
		return err
	}
	d.db = db

	return err
}

func (d *DataStore) Close() error {
	return d.db.Close()
}

// StoreTrade stores a trader to the database
func (d *DataStore) StoreTrade(t *Trade) *gorm.DB {
	return d.db.Create(t)
}

// CreateDB creates the database if it is not exist
func CreateDB() {
}

type currencies struct {
	Name       string
	Currencies []string
}

package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DataStore ...
type DataStore struct {
	Dialect  string
	Name     string
	User     string
	Password string

	db *gorm.DB
}

func NewDataStore(Dialet string) *DataStore {
	return &DataStore{Dialect: Dialet,
		Name:     "exdata",
		User:     "root",
		Password: "office98"} //@TODO later should change to from the config
}

func (d *DataStore) OpenDB() error {
	//	user+":"+passwd+"@tcp(127.0.0.1:3306)/"+name+"?charset=utf8&parseTime=True&loc=Local"
	uri := d.User + ":" + d.Password + "@/" + d.Name + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(d.Dialect, uri)
	if err != nil {
		log.Println("DB open with", err)
		return err
	}
	d.db = db

	return nil
}

func (d *DataStore) CloseDB() error {
	return d.db.Close()
}

func (d *DataStore) GetDB() *gorm.DB {
	return d.db
}

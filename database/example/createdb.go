package main

import (
	_ "github.com/go-sql-driver/mysql"
)

/*
func fake_main() {
	name := "exdata"
	passwd := "office98"
	user := "root"

	db, err := sql.Open("mysql", user+":"+passwd+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("Create Database If Not Exists " + name + "Character Set UTF8 ")
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
}*/

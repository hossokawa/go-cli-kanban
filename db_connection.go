package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", "kanban.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	exists, err := tableExists(db)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		if err := createTable(db); err != nil {
			log.Fatal(err)
		}
	}
}

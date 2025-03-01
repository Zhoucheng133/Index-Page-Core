package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitSql() {
	db, err := sql.Open("sqlite3", "db/pages.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	createPage := `CREATE TABLE IF NOT EXISTS pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		port TEXT,
		webui INTEGER,
		tip TEXT
	);`
	_, err = db.Exec(createPage)
	if err != nil {
		log.Fatal(err)
	}
	createUser := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		password TEXT,
	);`
	_, err = db.Exec(createUser)
	if err != nil {
		log.Fatal(err)
	}
}

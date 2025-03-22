package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitSql() {
	var err error
	db, err = sql.Open("sqlite3", "db/database.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	createPage := `CREATE TABLE IF NOT EXISTS pages (
		id TEXT PRIMARY KEY,
		icon TEXT,
		name TEXT,
		port TEXT,
		webui TEXT,
		tip TEXT
	);`
	_, err = db.Exec(createPage)
	if err != nil {
		log.Fatal(err)
	}
	createUser := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT,
		password TEXT,
		salt TEXT
	);`
	_, err = db.Exec(createUser)
	if err != nil {
		log.Fatal(err)
	}
}

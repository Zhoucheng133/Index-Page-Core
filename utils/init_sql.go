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
	}
	defer db.Close()
	createTable := `CREATE TABLE IF NOT EXISTS pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		port TEXT,
		webui INTEGER,
		tip TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

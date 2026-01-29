package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		uuid TEXT UNIQUE,
		password TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDetails(uuid string) (map[string]any, error) {
    row := DB.QueryRow(
        "SELECT name FROM users WHERE uuid = ?",
        uuid,
		) 

    var name string
    err := row.Scan(&name)
    if err != nil {
        return nil, err
    }

    return map[string]any{
        "Name": name,
    }, nil
}

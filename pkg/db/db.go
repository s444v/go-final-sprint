package db

import (
	"database/sql"
)

const scheme = `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title VARCHAR NOT NULL,
        comment TEXT,
        repeat VARCHAR CHECK(length(repeat) <= 128)
    );
    `

func dbConnect(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	_, err = db.Exec(scheme)
	if err != nil {
		return err
	}
	return err
}

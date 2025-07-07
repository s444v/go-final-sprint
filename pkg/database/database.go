package database

import (
	"database/sql"
)

const scheme = `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date CHAR(8) NOT NULL,
        title VARCHAR(255) NOT NULL,
        comment TEXT,
        repeat VARCHAR(128) NOT NULL CHECK(length(repeat) <= 128)
    );
    `

func DbConnect(dbFile string) error {
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

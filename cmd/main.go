package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/s444v/go-final-sprint/internal/server"
	_ "modernc.org/sqlite"
)

func main() {
	logger := log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Llongfile)
	err := dbConnect(dbFileName)
	if err != nil {
		logger.Fatalf("Ошибка в работе с базой данных: %v", err)
	}
	defer db.Close()
	server := server.NewServer(logger)
	err = server.HTTP.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

// -- db file
const scheme = `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT CHECK(length(repeat) <= 128)
    );
    `

var db *sql.DB
var dbFileName = os.Getenv("TODO_DBFILE")

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

// -- db file

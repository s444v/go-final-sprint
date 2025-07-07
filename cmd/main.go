package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/s444v/go-final-sprint/pkg/database"
	"github.com/s444v/go-final-sprint/pkg/server"
	_ "modernc.org/sqlite"
)

var db *sql.DB
var dbFileName = func() string {
	if f := os.Getenv("TODO_DBFILE"); f != "" {
		return f
	}
	return "scheduler.db"
}()

func main() {
	logger := log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Llongfile)
	err := database.DbConnect(dbFileName)
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

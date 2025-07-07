package server

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Logger *log.Logger
	HTTP   *http.Server
}

const WEBDIR = "./web"

var PORT = ":" + func() string {
	if p := os.Getenv("TODO_PORT"); p != "" {
		return p
	}
	return "7540"
}()

func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(WEBDIR)))

	httpServer := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		HTTP:   httpServer,
	}
}

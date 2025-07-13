package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/s444v/go-final-sprint/pkg/api"
)

type Server struct {
	Logger *log.Logger
	HTTP   *http.Server
}

var PORT = ":" + func() string {
	if p := os.Getenv("TODO_PORT"); p != "" {
		return p
	}
	return "7540"
}()

// Создаем сервер
func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	// Добавлем обработчики
	api.HandlersInit(mux)

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

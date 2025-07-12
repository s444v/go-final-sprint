package api

import (
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

// Обработчик для получения задачи по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	task, err := database.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// Возвращаем таску json файлом
	w.WriteHeader(http.StatusOK)
	writeJSON(w, task)
}

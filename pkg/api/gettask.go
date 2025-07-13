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
		err := writeJSON(w, map[string]string{"error": "id is required"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	task, err := database.GetTask(id)
	if err != nil {
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// Возвращаем таску json файлом
	w.WriteHeader(http.StatusOK)
	err = writeJSON(w, task)
	if err != nil {
		http.Error(w, "cant parse to json", http.StatusInternalServerError)
	}
}

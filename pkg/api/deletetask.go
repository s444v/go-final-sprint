package api

import (
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

// Обработчик для удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	// Вызов функции для удаления задачи по id
	if err := database.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// возвращаем пустой json если все хорошо
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{})
}

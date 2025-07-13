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
		err := writeJSON(w, map[string]string{"error": "id is required"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// Вызов функции для удаления задачи по id
	if err := database.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// возвращаем пустой json если все хорошо
	w.WriteHeader(http.StatusAccepted)
	err := writeJSON(w, map[string]string{})
	if err != nil {
		http.Error(w, "cant parse to json", http.StatusInternalServerError)
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

// Обработчик для редактирования задачи
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// если нет названия
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "Ошибка"})
		return
	}
	// проверка даты на валидность
	if err = checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// вызываем функция для обновления таски в базе данных
	err = database.UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// возвращаем пустой json если все хорошо
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{})
}

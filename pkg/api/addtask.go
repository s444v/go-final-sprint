package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

// Обработчик для добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	// Нет названия задачи
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "Ошибка"})
		return
	}
	// Проверка на валидность даты
	if err = checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// Добавление в базу данных
	id, err := database.AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// Возвращаем id добавленной задачи
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

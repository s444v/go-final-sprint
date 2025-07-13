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
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// если нет названия
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = writeJSON(w, map[string]string{"error": "Ошибка"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// проверка даты на валидность
	if err = checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// вызываем функция для обновления таски в базе данных
	err = database.UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// возвращаем пустой json если все хорошо
	w.WriteHeader(http.StatusAccepted)
	err = writeJSON(w, map[string]string{})
	if err != nil {
		http.Error(w, "cant parse to json", http.StatusInternalServerError)
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

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
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "Ошибка"})
		return
	}
	if err = checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	err = database.UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{})
}

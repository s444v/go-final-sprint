package api

import (
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	if err := database.DeleteTask(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{})
}

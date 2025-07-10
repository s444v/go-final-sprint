package api

import (
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

type TasksResp struct {
	Tasks []*database.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.Tasks(50, r.FormValue("search")) // в параметре максимальное количество записей
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

package api

import (
	"net/http"

	"github.com/s444v/go-final-sprint/pkg/database"
)

// Структура массива задач, для возврата списка задач
type TasksResp struct {
	Tasks []*database.Task `json:"tasks"`
}

// Обработчик для получения списка задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetTasks(50, r.FormValue("search")) // в параметре максимальное количество записей
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// Возвращаем список тасков json файлом
	w.WriteHeader(http.StatusOK)
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

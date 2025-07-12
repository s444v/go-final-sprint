package api

import (
	"net/http"
	"time"

	"github.com/s444v/go-final-sprint/pkg/database"
)

const TIMEFORMAT = "20060102"
const WEBDIR = "./web"

// Инициализация обработчиков
func HandlersInit(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(WEBDIR)))
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", getTasksHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
}

// Обработчик для поиска след. даты задачи
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(TIMEFORMAT, r.FormValue("now"))
	if err != nil {
		http.Error(w, "cant parse time", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if repeat == "" {
		w.WriteHeader(http.StatusOK) //тут исправить
		return
	}
	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "cant find next date", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// Распределитель по методам для запроса "/api/task"
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

// Обработчик для отметки о выполнении задачи
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	task, err := database.GetTask(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// Если у задачи нет заданного повторения = удаляем
	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusAccepted)
		writeJSON(w, map[string]string{})
		return
	}
	// Ищем след. дату для задачи
	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	// Обновляем базу данных с новой датой
	err = database.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{})
}

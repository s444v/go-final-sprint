package api

import (
	"net/http"
	"time"

	"github.com/s444v/go-final-sprint/pkg/database"
)

const TIMEFORMAT = "20060102"
const WEBDIR = "./web"

func HandlersInit(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(WEBDIR)))
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", getTasksHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(TIMEFORMAT, r.FormValue("now"))
	if err != nil {
		http.Error(w, "cant parse time", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if repeat == "" {
		w.WriteHeader(http.StatusOK)
	}
	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "cant find next date", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

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

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	task, err := database.GetTask(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
	}
	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}
	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	err = database.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{})
}

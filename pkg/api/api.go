package api

import (
	"net/http"
	"time"
)

const TIMEFORMAT = "20060102"
const WEBDIR = "./web"

func HandlersInit(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(WEBDIR)))
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
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
	}
}

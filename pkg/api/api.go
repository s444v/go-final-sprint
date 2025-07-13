package api

import (
	"net/http"
	"os"
	"time"

	"github.com/s444v/go-final-sprint/pkg/database"
)

const TIME_FORMAT = "20060102"
const WEB_DIR = "./web"
const JWT_SECRET = "12345"

var TODO_PASSWORD = func() string {
	if p := os.Getenv("TODO_PASSWORD"); p != "" {
		return p
	}
	return "12345"
}()

// Инициализация обработчиков
func HandlersInit(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(WEB_DIR)))
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/tasks", auth(getTasksHandler))
	mux.HandleFunc("/api/task/done", auth(doneTaskHandler))
	mux.HandleFunc("/api/signin", signinHandler)
}

// Обработчик для поиска след. даты задачи
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}
	now, err := time.Parse(TIME_FORMAT, r.FormValue("now"))
	if err != nil {
		http.Error(w, "cant parse time", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if repeat == "" {
		w.WriteHeader(http.StatusConflict) //тут исправить
		return
	}
	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "cant find next date", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(result))
	if err != nil {
		http.Error(w, "cant parse to json", http.StatusInternalServerError)
	}
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
	default:
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}
}

// Обработчик для отметки о выполнении задачи
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}
	id := r.FormValue("id")
	task, err := database.GetTask(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	// Если у задачи нет заданного повторения = удаляем
	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = writeJSON(w, map[string]string{"error": err.Error()})
			if err != nil {
				http.Error(w, "cant parse to json", http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusAccepted)
		err = writeJSON(w, map[string]string{})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
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
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	err = writeJSON(w, map[string]string{})
	if err != nil {
		http.Error(w, "cant parse to json", http.StatusInternalServerError)
	}
}

package api

import "net/http"

const timeFormat = "20060102"
const WEBDIR = "./web"

func HandlersInit(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(WEBDIR)))
	mux.HandleFunc("/api/nextdate", nextDayHandler)
}

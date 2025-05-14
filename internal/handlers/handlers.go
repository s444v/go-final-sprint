package handlers

import (
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.FileServer(http.Dir(webDir))
	// 	http.NotFound(w, r)
	// 	return
	// }
	// http.ServeFile(w, r, "index.html")
}

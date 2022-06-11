package handlers

import "net/http"

func HandlerStartPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("<h1>Welcome to gophermart</h1>"))
}

package server

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, World!\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page!\n"))
}

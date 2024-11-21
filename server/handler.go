package server

import (
	"html/template"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, World!\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./view/templates/index.html"))
	tmpl.Execute(w, nil)
}

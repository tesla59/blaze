package handler

import "net/http"

type ClientHandler interface {
	Handle() func(w http.ResponseWriter, r *http.Request)
}

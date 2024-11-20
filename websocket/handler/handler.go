package handler

import "net/http"

type WebsocketHandler interface {
	Handle() func(w http.ResponseWriter, r *http.Request)
}

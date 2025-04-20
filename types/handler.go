package types

import "net/http"

type Handler interface {
	Handle() func(w http.ResponseWriter, r *http.Request)
}

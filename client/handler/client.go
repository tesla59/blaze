package handler

import "net/http"

type ClientWrapper struct{}

func NewClientHandler() *ClientWrapper {
	return &ClientWrapper{}
}

func (c *ClientWrapper) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Client Handler"))
	}
}

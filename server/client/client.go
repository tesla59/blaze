package client

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

// Handler handles client-related operations.
type Handler struct {
	DB *pgxpool.Pool
}

// NewClientHandler creates a new ClientHandler with the provided database connection pool.
func NewClientHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		DB: db,
	}
}

func (c *Handler) Handle(method string) http.HandlerFunc {
	switch method {
	case "POST":
		return func(w http.ResponseWriter, r *http.Request) {
			c.postHandler(w, r)
		}
	case "GET":
		return func(w http.ResponseWriter, r *http.Request) {
			c.getHandler(w, r)
		}
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (c *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented yet\n"))
}

func (c *Handler) postHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented yet\n"))
}

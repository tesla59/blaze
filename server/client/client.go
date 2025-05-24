package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/repository"
	"github.com/tesla59/blaze/service"
	"log/slog"
	"net/http"
)

// Handler handles client-related operations.
type Handler struct {
	service *service.ClientService
	DB      *pgxpool.Pool
}

// NewClientHandler creates a new ClientHandler with the provided database connection pool.
func NewClientHandler(db *pgxpool.Pool) *Handler {
	repo := repository.NewClientRepository(db)
	service := service.NewClientService(repo)
	return &Handler{
		DB:      db,
		service: service,
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
	ctx := r.Context()
	client, err := c.service.RegisterClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client.Token = signUser(client.ID, client.UUID, client.UserName, "tesla")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(client)
	if err != nil {
		slog.Error("Error encoding response", "client", client, "error", err)
	}
}

func signUser(id int, uuid, username, secret string) string {
	payload := fmt.Sprintf("%d|%s|%s", id, uuid, username)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

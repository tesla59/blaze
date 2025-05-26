package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/models"
	"github.com/tesla59/blaze/repository"
	"github.com/tesla59/blaze/service"
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
			switch r.URL.Path {
			case "/api/v1/client":
				c.postHandler(w, r)
			case "/api/v1/client/verify":
				c.verifyHandler(w, r)
			}
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
	ctxLogger := log.Logger.With("method", "postHandler", "path", r.URL.Path)
	ctx := log.Inject(r.Context(), ctxLogger)

	log.WithContext(ctx).Info("Handling client registration request")

	client, err := c.service.RegisterClient(ctx)
	if err != nil {
		log.WithContext(ctx).Error("Error registering client", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client.Token = signUser(client.ID, client.UUID, client.UserName, config.GetConfig().Server.Secret)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(client)
	if err != nil {
		log.WithContext(ctx).Error("Error encoding response", "client", client, "error", err)
	}
	log.WithContext(ctx).Info("Client registered", "client", client)
}

func (c *Handler) verifyHandler(w http.ResponseWriter, r *http.Request) {
	ctxLogger := log.Logger.With("method", "verifyHandler", "path", r.URL.Path)
	ctx := log.Inject(r.Context(), ctxLogger)

	log.WithContext(ctx).Info("Handling client token verification request")

	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		log.WithContext(ctx).Error("Error decoding request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Verify the token
	expectedToken := signUser(client.ID, client.UUID, client.UserName, config.GetConfig().Server.Secret)
	if client.Token != expectedToken {
		log.WithContext(ctx).Warn("Invalid token", "clientID", client.ID, "expectedToken", expectedToken, "receivedToken", client.Token)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token verified successfully\n"))
	log.WithContext(ctx).Info("Token verified successfully", "clientID", client.ID)
}

func signUser(id int, uuid, username, secret string) string {
	payload := fmt.Sprintf("%d|%s|%s", id, uuid, username)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

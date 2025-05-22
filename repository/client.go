package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/models"
)

// ClientRepository defines the interface for client-related database operations.
type ClientRepository interface {
	Create(client *models.Client) error
}

// clientRepository implements the ClientRepository interface using a PostgreSQL database.
type clientRepository struct {
	db *pgxpool.Pool
}

// NewClientRepository creates a new instance of clientRepository with the provided database connection pool.
func NewClientRepository(db *pgxpool.Pool) ClientRepository {
	return &clientRepository{
		db: db,
	}
}

// Create inserts a new client into the database.
func (c clientRepository) Create(client *models.Client) error {
	return errors.New("not implemented")
}

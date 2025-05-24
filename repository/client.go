package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/models"
)

// ClientRepository defines the interface for client-related database operations.
type ClientRepository interface {
	Create(context.Context, *models.Client) error
	GetClientByID(context.Context, int) (*models.Client, error)
	GetClientByUUID(context.Context, string) (*models.Client, error)
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
func (c clientRepository) Create(ctx context.Context, client *models.Client) error {
	query := `INSERT INTO clients (uuid, username) VALUES ($1, $2) RETURNING id`
	err := c.db.QueryRow(ctx, query, client.UUID, client.UserName).Scan(&client.ID)
	return err
}

func (c clientRepository) GetClientByID(ctx context.Context, id int) (*models.Client, error) {
	query := `SELECT id, uuid, username FROM clients WHERE id = $1`
	client := &models.Client{}
	err := c.db.QueryRow(ctx, query, id).Scan(&client.ID, &client.UUID, &client.UserName)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c clientRepository) GetClientByUUID(ctx context.Context, uuid string) (*models.Client, error) {
	//TODO implement me
	panic("implement me")
}

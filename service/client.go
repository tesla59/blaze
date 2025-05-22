package service

import (
	"github.com/tesla59/blaze/models"
	"github.com/tesla59/blaze/repository"
)

// ClientService acts as a service layer for managing clients.
type ClientService struct {
	repo repository.ClientRepository
}

// NewClientService creates a new instance of ClientService.
func NewClientService(repo repository.ClientRepository) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

// RegisterClient registers a new client with the given ID, UUID, and username.
func (s *ClientService) RegisterClient(id, uuid, username string) error {
	return s.repo.Create(&models.Client{
		ID:       id,
		UUID:     uuid,
		UserName: username,
	})
}

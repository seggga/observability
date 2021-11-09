package service

import (
	"github.com/seggga/observability/internal/pkg/model"
)

type Storage interface {
	// CreateUser(*model.User) (*uuid.UUID, error) // creates a new user
	// DeleteUser(*uuid.UUID) error                // deletes a new user
	CreateLink(*model.Link) error        // creates a new redirect link
	GetLink(string) (*model.Link, error) // retrieves all data that corresponds to the short link
	DeleteLink(string) error             // deletes the link specified
	Close()                              // close connection to the storage (database / file / ....)
	IsSet(string) bool                   // checks if the short ID is in the database
}

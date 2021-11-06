package service

import (
	"github.com/seggga/observability/internal/pkg/model"
)

type Storage interface {
	// Init()
	// CreateUser(*model.User) (*uuid.UUID, error) // creates a new user
	// DeleteUser(*uuid.UUID) error                // deletes a new user
	CreateLink(*model.Link) error        // creates a new redirect link
	GetLink(string) (*model.Link, error) // retrieves all data that corresponds to the short link
	//	Redirect(string) (string, error)              // retrieves long URL from database to produce redirect
	// DeleteLink(string, *uuid.UUID) error // deletes the link specified
	Close()            // close connection to the storage (database / file / ....)
	IsSet(string) bool // checks if the short ID is in the database
}

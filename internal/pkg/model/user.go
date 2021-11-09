package model

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID `json:"ID"`
	Name string    `json:"Name"`
}

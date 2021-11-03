package model

import "github.com/gofrs/uuid"

type Link struct {
	Short       string    `json:"Short"`
	Long        string    `json:"Long"`
	Owner       uuid.UUID `json:"Owner"`
	Count       int64     `json:"Count"`
	Description string    `json:"Description"`
}

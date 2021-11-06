package cropper

import "github.com/google/uuid"

type Link struct {
	Short       string    `json:"short"`
	Long        string    `json:"long"`
	Owner       uuid.UUID `json:"owner"`
	Count       int64     `json:"count"`
	Description string    `json:"description"`
}

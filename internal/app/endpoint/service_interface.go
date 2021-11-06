package endpoint

import "github.com/seggga/observability/pkg/cropper"

type service interface {
	Resolve(string) (string, error)
	NewLink(*cropper.Link) error
}

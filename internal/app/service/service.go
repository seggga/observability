package service

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/seggga/observability/internal/pkg/model"
	"github.com/seggga/observability/pkg/cropper"
)

type Service struct {
	repo Storage // repository holding application data
}

// New creates a new variable of type Service
func New(repo Storage) *Service {
	return &Service{
		repo: repo,
	}
}

// Resolve retrieves long-link from the storage layer on given short-link
func (s *Service) Resolve(short string) (string, error) {
	link, err := s.repo.GetLink(short)
	if err != nil {
		return "", nil
	}
	return link.Long, nil
}

// NewLink checks given data (short and long URIs) and sends data to the storage layer
func (s *Service) NewLink(link *cropper.Link) error {

	long, err := url.Parse(link.Long)
	// check user's input: incorrect URL format
	if err != nil {
		err = fmt.Errorf("entered long URL cannot be recognized: (%s) %w", link.Long, err)
		// slogger.Debug(err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}
	// check user's input: schema is empty
	if long.Scheme == "" {
		err := fmt.Errorf("protocol should be set (http:// or https:// or ...): %s", link.Long)
		// slogger.Debug(err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}
	// check user's input: host not set
	if long.Host == "" {
		err := fmt.Errorf("host address was not set: %s", link.Long)
		// slogger.Debug(err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}

	// check if short is not in use
	if s.repo.IsSet(link.Short) {
		err := fmt.Errorf("short URL %s is already in use", link.Short)
		// slogger.Debug(err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}

	// write data to the storage
	if err := s.repo.CreateLink(&model.Link{
		Short:       link.Short,
		Long:        link.Long,
		Owner:       link.Owner,
		Count:       link.Count,
		Description: link.Description,
	}); err != nil {
		err = fmt.Errorf("error creating new short-to-long pair %w", err)
		// slogger.Errorw("error creating new short-to-long pair", err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}
	return nil
}

// DeleteLink deletes link from the storage
func (s *Service) DeleteLink(short string, ID *uuid.UUID) error {
	// check link existance
	if !s.repo.IsSet(short) {
		err := fmt.Errorf("short URL you wish to delete (%s) has not been found in the storage. Nothing to delete", short)
		// slogger.Debug(err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}

	// call link delete
	err := s.repo.DeleteLink(short)
	return err
}

package service

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/seggga/observability/internal/pkg/model"
	"github.com/seggga/observability/pkg/cropper"
	"go.uber.org/zap"
)

type Service struct {
	repo   Storage // repository holding application data
	logger *zap.Logger
}

// New creates a new variable of type Service
func New(repo Storage, logLevel string) *Service {
	logger := initLogger(logLevel)
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Resolve retrieves long-link from the storage layer on given short-link
func (s *Service) Resolve(short string) (string, error) {

	s.logger.Sugar().Debugf("resolve %s called", short)
	link, err := s.repo.GetLink(short)
	if err != nil {
		s.logger.Sugar().Errorf("error resolving short link: %v", err)
		return "", nil
	}
	s.logger.Sugar().Debugf("found link %v", link)
	return link.Long, nil
}

// NewLink checks given data (short and long URIs) and sends data to the storage layer
func (s *Service) NewLink(link *cropper.Link) error {

	long, err := url.Parse(link.Long)
	if err != nil {
		err = fmt.Errorf("entered long URL cannot be recognized: (%s) %w", link.Long, err)
		s.logger.Error(err.Error())
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}
	// check user's input: schema is empty
	if long.Scheme == "" {
		err := fmt.Errorf("protocol should be set (http:// or https:// or ...): %s", link.Long)
		s.logger.Error(err.Error())
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}
	// check user's input: host not set
	if long.Host == "" {
		err := fmt.Errorf("host address was not set: %s", link.Long)
		s.logger.Error(err.Error())
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}

	// check if short is not in use
	if s.repo.IsSet(link.Short) {
		err := fmt.Errorf("short URL %s is already in use", link.Short)
		s.logger.Error(err.Error())
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}

	// write data to the storage
	if err := s.repo.CreateLink(
		&model.Link{
			Short:       link.Short,
			Long:        link.Long,
			Owner:       link.Owner,
			Count:       link.Count,
			Description: link.Description,
		},
	); err != nil {
		err = fmt.Errorf("error creating new short-to-long pair %w", err)
		s.logger.Error(err.Error())
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}
	return nil
}

// DeleteLink deletes link from the storage
func (s *Service) DeleteLink(short string, ID *uuid.UUID) error {
	s.logger.Sugar().Debugf("delete %s link called", short)
	// check link existance
	if !s.repo.IsSet(short) {
		err := fmt.Errorf("short URL you wish to delete (%s) has not been found in the storage. Nothing to delete", short)
		s.logger.Sugar().Error(err)
		// JSONError(rw, err, http.StatusBadRequest)
		return err
	}

	// call link delete
	err := s.repo.DeleteLink(short)
	if err != nil {
		s.logger.Error(err.Error())
	}
	return err
}

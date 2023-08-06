package shorten

import (
	"context"
	"github.com/google/uuid"
	"github.com/winterochek/go-linker/internal/model"
)

type Storage interface {
	Put(ctx context.Context, shortening model.Shortening) (*model.Shortening, error)
	Get(ctx context.Context, identifier string) (*model.Shortening, error)
	IncrementVisits(ctx context.Context, identifier string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(ctx context.Context, input model.ShortenInput) (*model.Shortening, error) {
	// generate identifier
	var (
		id         = uuid.New().ID()
		identifier = input.Identifier.OrElse(Shorten(id))
	)
	// save in storage
	dbShortening := model.Shortening{
		Identifier:  identifier,
		OriginalURL: input.RawURL,
		CreatedBy:   input.CreatedBy,
	}
	shortening, err := s.storage.Put(ctx, dbShortening)
	if err != nil {
		return nil, err
	}
	return shortening, nil
}

func (s *Service) Redirect(ctx context.Context, identifier string) (string, error) {
	// Get original URL
	shortening, err := s.storage.Get(ctx, identifier)
	if err != nil {
		return "", err
	}

	// increase views count
	if err := s.storage.IncrementVisits(ctx, identifier); err != nil {
		return "", err
	}

	return shortening.OriginalURL, nil

}

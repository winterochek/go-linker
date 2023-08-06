package shorten

import (
	"context"
	"github.com/winterochek/go-linker/internal/model"
	"github.com/google/uuid"
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
	}
	shortening, err := s.storage.Put(ctx, dbShortening)
	if err != nil {
		return nil, err
	}
	return shortening, nil
}

package product

import (
	"context"
	"errors"

	"sampleBackend/internal/storage"
)

var (
	ErrExist = errors.New("item exist")
)

type Storage interface {
	Create(ctx context.Context, p Product) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) AddProduct(ctx context.Context, p Product) error {
	err := s.storage.Create(ctx, p)
	if err != nil {
		if storage.IsErrAlreadyExist(err) {
			return ErrExist
		}
		return err
	}

	return nil
}

func IsErrExist(err error) bool {
	return errors.Is(err, ErrExist)
}

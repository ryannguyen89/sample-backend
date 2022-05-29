package user

import (
	"context"
	"errors"
	"fmt"

	"sampleBackend/internal/storage"
)

var (
	ErrUserExist = errors.New("user exist")
)

type Storage interface {
	Create(ctx context.Context, u User) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) CreateUser(ctx context.Context, u User) error {
	err := s.storage.Create(ctx, u)
	if err != nil {
		if storage.IsErrAlreadyExist(err) {
			return fmt.Errorf("create user: %v - %w", err, ErrUserExist)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func IsErrUserExist(err error) bool {
	return errors.Is(err, ErrUserExist)
}

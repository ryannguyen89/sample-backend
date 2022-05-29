package product

import (
	"context"
	"errors"

	"sampleBackend/internal/storage"
)

var (
	ErrExist    = errors.New("item exist")
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	Create(ctx context.Context, p Product) error
	Update(ctx context.Context, p Product) error
	Delete(ctx context.Context, sku string) error
	List(ctx context.Context) ([]*Product, error)
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

func (s *Service) UpdateProduct(ctx context.Context, p Product) error {
	err := s.storage.Update(ctx, p)
	if err != nil {
		if storage.IsErrNotFound(err) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, sku string) error {
	err := s.storage.Delete(ctx, sku)
	if err != nil {
		if storage.IsErrNotFound(err) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *Service) ListProduct(ctx context.Context) ([]*Product, error) {
	return s.storage.List(ctx)
}

func IsErrExist(err error) bool {
	return errors.Is(err, ErrExist)
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

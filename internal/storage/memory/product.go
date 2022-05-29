package memory

import (
	"context"
	"sync"

	"sampleBackend/internal/product"
	"sampleBackend/internal/storage"
)

type ProductStorage struct {
	mu sync.Mutex

	products map[string]product.Product
}

func NewProductStorage() *ProductStorage {
	return &ProductStorage{
		products: make(map[string]product.Product),
	}
}

func (ps *ProductStorage) Create(_ context.Context, p product.Product) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exist := ps.products[p.SKU]; exist {
		return storage.ErrAlreadyExist
	}

	ps.products[p.SKU] = p
	return nil
}

func (ps *ProductStorage) Update(_ context.Context, p product.Product) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exist := ps.products[p.SKU]; !exist {
		return storage.ErrNotFound
	}

	ps.products[p.SKU] = p
	return nil
}

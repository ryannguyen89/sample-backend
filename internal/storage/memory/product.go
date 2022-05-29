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

func (ps *ProductStorage) Get(_ context.Context, sku string) (*product.Product, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if item, exist := ps.products[sku]; exist {
		item := item
		return &item, nil
	} else {
		return nil, storage.ErrNotFound
	}
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

func (ps *ProductStorage) Delete(_ context.Context, sku string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exist := ps.products[sku]; !exist {
		return storage.ErrNotFound
	}

	delete(ps.products, sku)

	return nil
}

func (ps *ProductStorage) List(_ context.Context) ([]*product.Product, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var retList []*product.Product
	for _, p := range ps.products {
		pTemp := p
		retList = append(retList, &pTemp)
	}

	return retList, nil
}

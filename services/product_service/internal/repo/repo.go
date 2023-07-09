package repo

import (
	"context"
	"sync"

	"salespot/services/product_service/internal/models"
	"salespot/shared/sctx/component/tracing"
)

type DbStorage interface {
	ListProduct(ctx context.Context) ([]models.Product, error)
	GetProduct(ctx context.Context, id string) (*models.Product, error)
}

type CacheStorage interface {
	GetProduct(ctx context.Context, id string) (*models.Product, error)
	StoreProduct(ctx context.Context, product *models.Product) error
}

type repository struct {
	store        DbStorage
	cacheStorage CacheStorage
	once         sync.Once
}

func NewRepository(store DbStorage, cacheStorage CacheStorage) *repository {
	return &repository{store: store, cacheStorage: cacheStorage}
}

func (r *repository) ListProduct(ctx context.Context) ([]models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "repo.list-product")
	defer span.End()

	return r.store.ListProduct(ctx)
}

func (r *repository) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "repo.get-product")
	defer span.End()

	shouldCache := false
	product, err := r.cacheStorage.GetProduct(ctx, id)
	if err != nil {
		shouldCache = true
		product, err = r.store.GetProduct(ctx, id)
	}

	if shouldCache {
		r.once.Do(func() {
			_ = r.cacheStorage.StoreProduct(ctx, product)
		})
	}

	return product, nil
}

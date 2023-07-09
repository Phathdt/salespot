package repo

import (
	"context"

	"salespot/services/product_service/internal/models"
)

type MongoStorage interface {
	ListProduct(ctx context.Context) ([]models.Product, error)
	GetProduct(ctx context.Context, id string) (*models.Product, error)
}

type repository struct {
	store MongoStorage
}

func NewRepository(store MongoStorage) *repository {
	return &repository{store: store}
}

func (r *repository) ListProduct(ctx context.Context) ([]models.Product, error) {
	return r.store.ListProduct(ctx)
}

func (r *repository) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	return r.store.GetProduct(ctx, id)
}

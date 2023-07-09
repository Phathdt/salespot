package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"salespot/services/product_service/internal/models"
	"salespot/shared/sctx/component/tracing"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *redisStore {
	return &redisStore{client: client}
}

func (r *redisStore) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "cache-storage.get-product")
	defer span.End()

	result, err := r.client.Get(ctx, fmt.Sprintf("/products/%s", id)).Result()
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err = bson.Unmarshal([]byte(result), &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *redisStore) StoreProduct(ctx context.Context, product *models.Product) error {
	ctx, span := tracing.StartTrace(ctx, "cache-storage.store-product")
	defer span.End()

	bytes, err := bson.Marshal(product)
	if err != nil {
		return err
	}

	if err = r.client.Set(ctx, fmt.Sprintf("/products/%s", product.ID.Hex()), bytes, time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

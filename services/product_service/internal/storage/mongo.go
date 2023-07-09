package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"salespot/services/product_service/internal/models"
)

type mongoStore struct {
	db *mongo.Database
}

func NewMongoStore(db *mongo.Database) *mongoStore {
	return &mongoStore{db: db}
}

func (m *mongoStore) ListProduct(ctx context.Context) ([]models.Product, error) {
	collection := m.db.Collection(models.Product{}.Collection())

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	products := make([]models.Product, 0)

	for cur.Next(ctx) {
		var product models.Product
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

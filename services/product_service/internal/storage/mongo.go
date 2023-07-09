package storage

import (
	"context"

	"salespot/services/product_service/internal/models"
	"salespot/shared/sctx/component/tracing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStore struct {
	db *mongo.Database
}

func NewMongoStore(db *mongo.Database) *mongoStore {
	return &mongoStore{db: db}
}

func (m *mongoStore) ListProduct(ctx context.Context) ([]models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "storage.list-product")

	defer span.End()

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

func (m *mongoStore) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "storage.get-product")
	defer span.End()

	collection := m.db.Collection(models.Product{}.Collection())

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

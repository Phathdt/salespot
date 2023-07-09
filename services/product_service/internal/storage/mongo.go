package storage

import (
	"context"

	"github.com/qiniu/qmgo"
	"salespot/services/product_service/internal/models"
	"salespot/shared/sctx/component/tracing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoStore struct {
	db *qmgo.Database
}

func NewMongoStore(db *qmgo.Database) *mongoStore {
	return &mongoStore{db: db}
}

func (m *mongoStore) ListProduct(ctx context.Context) ([]models.Product, error) {
	ctx, span := tracing.StartTrace(ctx, "storage.list-product")
	defer span.End()

	collection := m.db.Collection(models.Product{}.Collection())

	var products []models.Product
	if err := collection.Find(ctx, bson.M{}).All(&products); err != nil {
		return nil, err
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
	if err = collection.Find(ctx, bson.M{"_id": objectId}).One(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

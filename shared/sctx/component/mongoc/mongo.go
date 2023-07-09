package mongoc

import (
	"context"
	"flag"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"salespot/shared/sctx"
)

type MongoComponent interface {
	GetDb() *mongo.Database
}
type mongoDB struct {
	id            string
	prefix        string
	logger        sctx.Logger
	db            *mongo.Database
	client        *mongo.Client
	mongoURI      string
	mongoDatabase string
}

func NewMongoDB(id string, prefix string) *mongoDB {
	return &mongoDB{id: id, prefix: prefix}
}

func (m *mongoDB) ID() string {
	return m.id
}

func (m *mongoDB) InitFlags() {
	flag.StringVar(&m.mongoURI, "mongo_uri", "", "mongo uri")
	flag.StringVar(&m.mongoDatabase, "mongo_database", "", "mongo database")
}

func (m *mongoDB) Activate(sc sctx.ServiceContext) error {
	m.logger = sctx.GlobalLogger().GetLogger(m.id)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.mongoURI))
	if err != nil {
		fmt.Println(err)
		return err
	}

	m.client = client

	m.db = m.client.Database(m.mongoDatabase)

	return nil
}

func (m *mongoDB) Stop() error {
	return m.client.Disconnect(context.Background())
}

func (m *mongoDB) GetDb() *mongo.Database {
	return m.db
}

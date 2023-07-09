package mongoc

import (
	"context"
	"flag"
	"fmt"

	"github.com/qiniu/qmgo"
	qoptions "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"salespot/shared/sctx"
)

type MongoComponent interface {
	GetDb() *qmgo.Database
}
type mongoDB struct {
	id            string
	prefix        string
	logger        sctx.Logger
	db            *qmgo.Database
	client        *qmgo.Client
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

	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	clientOptions := qoptions.ClientOptions{ClientOptions: opts}
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: m.mongoURI}, clientOptions)
	if err != nil {
		fmt.Println(err)
		return err
	}

	m.client = client

	m.db = m.client.Database(m.mongoDatabase)

	return nil
}

func (m *mongoDB) Stop() error {
	return m.client.Close(context.Background())
}

func (m *mongoDB) GetDb() *qmgo.Database {
	return m.db
}

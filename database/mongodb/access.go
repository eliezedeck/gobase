package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/eliezedeck/gobase/database/mongodb/codecs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrFailedCreatingDBClient = errors.New("failed creating a DB client")
	ErrFailedConnecting       = errors.New("failed connecting to the DB")
)

type Access struct {
	dbClient *mongo.Client
	dbName   string
}

// NewDBAccess creates a new DBAccess, connect to the database
func NewDBAccess(uri, database string) (*Access, error) {
	builder := bson.NewRegistryBuilder()

	// Decoder for C# field value {"_csharpnull": true} to map to Go's *time.Time = nil
	builder.RegisterTypeDecoder(reflect.TypeOf(&time.Time{}), codecs.CSharpNullTimeDecoder{})

	clientOpts := options.Client().ApplyURI(uri).SetRegistry(builder.Build())

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, ErrFailedCreatingDBClient
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, ErrFailedConnecting
	}

	access := &Access{
		dbClient: client,
		dbName:   database,
	}
	return access, nil
}

func (a *Access) Collection(name string) *mongo.Collection {
	return a.dbClient.Database(a.dbName).Collection(name)
}

func (a *Access) CreateIndexes(ctx context.Context, coll string, models []mongo.IndexModel) error {
	c := a.Collection(coll)
	if _, err := c.Indexes().CreateMany(ctx, models); err != nil {
		return err
	}
	return nil
}

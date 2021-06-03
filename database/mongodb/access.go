package api

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DatabaseName = "bowa"

var (
	ErrFailedCreatingDBClient = errors.New("failed creating a DB client")
	ErrFailedConnecting       = errors.New("failed connecting to the DB")
)

type Access struct {
	dbClient *mongo.Client
}

// NewDBAccess creates a new DBAccess, connect to the database
func NewDBAccess(uri string) (*Access, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
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
	}
	return access, nil
}

func (a *Access) Collection(name string) *mongo.Collection {
	return a.dbClient.Database(DatabaseName).Collection(name)
}

func (a *Access) CreateIndexes(ctx context.Context, coll string, models []mongo.IndexModel) error {
	c := a.Collection(coll)
	if _, err := c.Indexes().CreateMany(ctx, models); err != nil {
		return err
	}
	return nil
}

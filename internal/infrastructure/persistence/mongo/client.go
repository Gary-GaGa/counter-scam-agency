package mongo

import (
	"context"
	"errors"
	"fmt"

	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config provides Mongo connection settings.
type Config struct {
	URI      string
	Database string
}

// NewClient creates and pings a Mongo client.
func NewClient(ctx context.Context, cfg Config) (*mongodriver.Client, error) {
	if cfg.URI == "" {
		return nil, errors.New("mongo uri is required")
	}
	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("ping mongo: %w", err)
	}
	return client, nil
}

// NewDatabase returns the configured database handle.
func NewDatabase(client *mongodriver.Client, cfg Config) (*mongodriver.Database, error) {
	if client == nil {
		return nil, errors.New("mongo client is nil")
	}
	if cfg.Database == "" {
		return nil, errors.New("mongo database is required")
	}
	return client.Database(cfg.Database), nil
}

// Disconnect closes the Mongo client.
func Disconnect(ctx context.Context, client *mongodriver.Client) error {
	if client == nil {
		return nil
	}
	return client.Disconnect(ctx)
}

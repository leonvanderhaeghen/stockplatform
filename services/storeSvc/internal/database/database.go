package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/config"
)

// Database wraps MongoDB client and provides database operations
type Database struct {
	client *mongo.Client
	db     *mongo.Database
}

// Initialize creates a new database connection
func Initialize(cfg *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.Database.Database)

	return &Database{
		client: client,
		db:     db,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

// GetCollection returns a MongoDB collection
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

// GetDatabase returns the MongoDB database instance
func (d *Database) GetDatabase() *mongo.Database {
	return d.db
}

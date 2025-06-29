package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/infrastructure/mongodb"
)

// Database holds database connections and repositories
type Database struct {
	Client          *mongo.Client
	Database        *mongo.Database
	ProductRepo     *mongodb.ProductRepository
	CategoryRepo    domain.CategoryRepository
	logger          *zap.Logger
}

// Initialize creates and initializes the database layer
func Initialize(cfg *config.Config, logger *zap.Logger) (*Database, error) {
	// Create MongoDB client
	client, err := createMongoClient(cfg.MongoURI, logger)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to MongoDB", zap.String("database", cfg.Database))

	// Get database instance
	database := client.Database(cfg.Database)

	// Initialize repositories
	productRepo := mongodb.NewProductRepository(database, logger)
	categoryRepo := mongodb.NewCategoryRepository(database, logger)

	return &Database{
		Client:       client,
		Database:     database,
		ProductRepo:  productRepo,
		CategoryRepo: categoryRepo,
		logger:       logger,
	}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Client.Disconnect(ctx); err != nil {
		db.logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return err
	}

	db.logger.Info("Disconnected from MongoDB")
	return nil
}

// createMongoClient creates a new MongoDB client with proper configuration
func createMongoClient(mongoURI string, logger *zap.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(10).
		SetMinPoolSize(1)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Failed to create MongoDB client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

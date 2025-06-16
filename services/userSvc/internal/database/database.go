package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
	mongorepo "github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/infrastructure/mongodb"
)

// Database holds database connections and repositories
type Database struct {
	Client         *mongo.Client
	Database       *mongo.Database
	UserRepo       domain.UserRepository
	AddressRepo    domain.AddressRepository
	PermissionRepo domain.PermissionRepository
	logger         *zap.Logger
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

	// Log MongoDB server status
	serverStatus, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to list MongoDB databases", zap.Error(err))
	} else {
		logger.Info("Successfully connected to MongoDB server",
			zap.Strings("databases", serverStatus),
		)
	}

	logger.Info("Successfully connected to MongoDB", zap.String("database", cfg.Database))

	// Get database instance
	database := client.Database(cfg.Database)

	// Initialize repositories
	userRepo := mongorepo.NewUserRepository(database, "users", logger)
	addressRepo := mongorepo.NewAddressRepository(database, "addresses", logger)
	permissionRepo := mongorepo.NewPermissionRepository(database, "permissions")

	return &Database{
		Client:         client,
		Database:       database,
		UserRepo:       userRepo,
		AddressRepo:    addressRepo,
		PermissionRepo: permissionRepo,
		logger:         logger,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Info("Disconnecting from MongoDB...")
	if err := d.Client.Disconnect(ctx); err != nil {
		d.logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return err
	}

	d.logger.Info("Successfully disconnected from MongoDB")
	return nil
}

// createMongoClient creates a MongoDB client with proper configuration
func createMongoClient(mongoURI string, logger *zap.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Connecting to MongoDB...")

	// Set up client options with timeouts and retry settings
	clientOptions := options.Client()
	clientOptions.ApplyURI(mongoURI)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)
	clientOptions.SetMaxPoolSize(10)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Failed to create MongoDB client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

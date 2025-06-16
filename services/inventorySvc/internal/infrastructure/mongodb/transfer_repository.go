package mongodb

import (
	"context"
	"errors"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// TransferRepository implements the domain.TransferRepository interface using MongoDB
type TransferRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewTransferRepository creates a new TransferRepository
func NewTransferRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.TransferRepository {
	return &TransferRepository{
		collection: db.Collection(collectionName),
		logger:     logger.Named("transfer_repository"),
	}
}

// Create inserts a new transfer in the database
func (r *TransferRepository) Create(ctx context.Context, transfer *domain.Transfer) error {
	r.logger.Debug("Creating transfer", zap.String("id", transfer.ID))

	if transfer.ID == "" {
		transfer.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, transfer)
	if err != nil {
		r.logger.Error("Failed to create transfer", zap.Error(err))
		return err
	}

	return nil
}

// GetByID retrieves a transfer by its ID
func (r *TransferRepository) GetByID(ctx context.Context, id string) (*domain.Transfer, error) {
	r.logger.Debug("Getting transfer by ID", zap.String("id", id))

	var transfer domain.Transfer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transfer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Transfer not found", zap.String("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get transfer", zap.Error(err))
		return nil, err
	}

	return &transfer, nil
}

// Update updates an existing transfer
func (r *TransferRepository) Update(ctx context.Context, transfer *domain.Transfer) error {
	r.logger.Debug("Updating transfer", zap.String("id", transfer.ID))

	filter := bson.M{"_id": transfer.ID}
	update := bson.M{"$set": transfer}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update transfer", zap.Error(err))
		return err
	}

	return nil
}

// Delete removes a transfer by its ID
func (r *TransferRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting transfer", zap.String("id", id))

	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete transfer", zap.Error(err))
		return err
	}

	return nil
}

// ListByStatus retrieves transfers by status
func (r *TransferRepository) ListByStatus(ctx context.Context, status domain.TransferStatus, limit, offset int) ([]*domain.Transfer, error) {
	r.logger.Debug("Listing transfers by status",
		zap.String("status", string(status)),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	filter := bson.M{"status": status}
	return r.listTransfers(ctx, filter, limit, offset)
}

// ListPendingTransfers retrieves all pending transfers
func (r *TransferRepository) ListPendingTransfers(ctx context.Context, limit, offset int) ([]*domain.Transfer, error) {
	r.logger.Debug("Listing pending transfers",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return r.ListByStatus(ctx, domain.TransferStatusRequested, limit, offset)
}

// ListBySourceLocation retrieves transfers by source location
func (r *TransferRepository) ListBySourceLocation(ctx context.Context, locationID string, limit, offset int) ([]*domain.Transfer, error) {
	r.logger.Debug("Listing transfers by source location",
		zap.String("location_id", locationID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	filter := bson.M{"source_location_id": locationID}
	return r.listTransfers(ctx, filter, limit, offset)
}

// ListByDestLocation retrieves transfers by destination location
func (r *TransferRepository) ListByDestLocation(ctx context.Context, locationID string, limit, offset int) ([]*domain.Transfer, error) {
	r.logger.Debug("Listing transfers by destination location",
		zap.String("location_id", locationID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	filter := bson.M{"destination_location_id": locationID}
	return r.listTransfers(ctx, filter, limit, offset)
}

// ListByProduct retrieves transfers by product ID
func (r *TransferRepository) ListByProduct(ctx context.Context, productID string, limit, offset int) ([]*domain.Transfer, error) {
	r.logger.Debug("Listing transfers by product",
		zap.String("product_id", productID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	filter := bson.M{"items.product_id": productID}
	return r.listTransfers(ctx, filter, limit, offset)
}

// Helper method to list transfers with a filter
func (r *TransferRepository) listTransfers(ctx context.Context, filter bson.M, limit, offset int) ([]*domain.Transfer, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.M{"requested_at": -1}) // Sort by requested date desc

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find transfers", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	transfers := make([]*domain.Transfer, 0)
	for cursor.Next(ctx) {
		var transfer domain.Transfer
		if err := cursor.Decode(&transfer); err != nil {
			r.logger.Error("Failed to decode transfer", zap.Error(err))
			continue
		}
		transfers = append(transfers, &transfer)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error when listing transfers", zap.Error(err))
		return nil, err
	}

	return transfers, nil
}

// GenerateID generates a new transfer ID
func (r *TransferRepository) GenerateID() string {
	return primitive.NewObjectID().Hex()
}

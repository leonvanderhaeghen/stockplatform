package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents a product category in the system
type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required,min=3,max=100"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	ParentID    string             `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Level       int32              `bson:"level" json:"level"`
	Path        string             `bson:"path" json:"path"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, parentID string, depth int32) ([]*Category, error)
}

// CategoryUseCase defines the business logic for category operations
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	GetCategory(ctx context.Context, id string) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context, parentID string, depth int32) ([]*Category, error)
}

package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Supplier represents a supplier in the system
type Supplier struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	ContactPerson string             `bson:"contact_person,omitempty" json:"contact_person,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty" validate:"omitempty,email"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Address       string             `bson:"address,omitempty" json:"address,omitempty"`
	City          string             `bson:"city,omitempty" json:"city,omitempty"`
	State         string             `bson:"state,omitempty" json:"state,omitempty"`
	Country       string             `bson:"country,omitempty" json:"country,omitempty"`
	PostalCode    string             `bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	TaxID         string             `bson:"tax_id,omitempty" json:"tax_id,omitempty"`
	Website       string             `bson:"website,omitempty" json:"website,omitempty"`
	Currency      string             `bson:"currency,omitempty" json:"currency,omitempty"`
	LeadTimeDays  int32              `bson:"lead_time_days,omitempty" json:"lead_time_days,omitempty"`
	PaymentTerms  string             `bson:"payment_terms,omitempty" json:"payment_terms,omitempty"`
	Metadata      map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// SupplierRepository defines the interface for supplier data operations
type SupplierRepository interface {
	// Create creates a new supplier
	Create(ctx context.Context, supplier *Supplier) (*Supplier, error)
	// GetByID retrieves a supplier by ID
	GetByID(ctx context.Context, id string) (*Supplier, error)
	// Update updates an existing supplier
	Update(ctx context.Context, supplier *Supplier) error
	// Delete deletes a supplier by ID
	Delete(ctx context.Context, id string) error
	// List retrieves a list of suppliers with pagination and optional filtering
	List(ctx context.Context, page, pageSize int32, search string) ([]*Supplier, int32, error)
}

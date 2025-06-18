package domain

// ResourceType represents the type of resource a permission applies to
type ResourceType string

const (
	// ResourceTypeSupplier represents a supplier resource
	ResourceTypeSupplier ResourceType = "SUPPLIER"
	// ResourceTypeProduct represents a product resource
	ResourceTypeProduct ResourceType = "PRODUCT"
)

// UserPermission represents a user's permission on a specific resource
type UserPermission struct {
	ID           string       `bson:"_id,omitempty"`
	UserID       string       `bson:"user_id"`
	ResourceType ResourceType `bson:"resource_type"`
	ResourceID   string       `bson:"resource_id"`
	Permission   string       `bson:"permission"`
}

// HasPermission checks if the permission includes the required permission level
func (p *UserPermission) HasPermission(required string) bool {
	// If the user has ADMIN permission, they have all permissions
	if p.Permission == "ADMIN" {
		return true
	}

	switch required {
	case "VIEW":
		return p.Permission == "VIEW" || 
		       p.Permission == "EDIT" || 
		       p.Permission == "DELETE"
	case "EDIT", "DELETE":
		return p.Permission == required
	default:
		return false
	}
}

// UserResource represents a resource that a user has access to
type UserResource struct {
	ID           string       `bson:"_id,omitempty"`
	UserID       string       `bson:"user_id"`
	ResourceType ResourceType `bson:"resource_type"`
	ResourceID   string       `bson:"resource_id"`
	Permission   string       `bson:"permission"`
}

// ResourceAccess represents a resource with the user's permission level
type ResourceAccess struct {
	ResourceID string     `json:"resource_id"`
	Permission string     `json:"permission"`
}

package domain

// Permission represents a user's permission on a resource
type Permission string

const (
	// PermissionView allows viewing a resource
	PermissionView Permission = "VIEW"
	// PermissionEdit allows editing a resource
	PermissionEdit Permission = "EDIT"
	// PermissionDelete allows deleting a resource
	PermissionDelete Permission = "DELETE"
	// PermissionAdmin gives full admin rights on a resource
	PermissionAdmin Permission = "ADMIN"
)

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
	Permission   Permission   `bson:"permission"`
}

// HasPermission checks if the permission includes the required permission level
func (p *UserPermission) HasPermission(required Permission) bool {
	// If the user has ADMIN permission, they have all permissions
	if p.Permission == PermissionAdmin {
		return true
	}

	switch required {
	case PermissionView:
		return p.Permission == PermissionView || 
		       p.Permission == PermissionEdit || 
		       p.Permission == PermissionDelete
	case PermissionEdit, PermissionDelete:
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
	Permission   Permission   `bson:"permission"`
}

// ResourceAccess represents a resource with the user's permission level
type ResourceAccess struct {
	ResourceID string     `json:"resource_id"`
	Permission Permission `json:"permission"`
}

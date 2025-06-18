package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Role represents user roles
type Role string

const (
	// RoleCustomer is a standard customer role
	RoleCustomer Role = "CUSTOMER"
	// RoleAdmin is an administrative role
	RoleAdmin Role = "ADMIN"
	// RoleStaff is a staff role with limited admin capabilities
	RoleStaff Role = "STAFF"
)

// User represents a user account
type User struct {
	ID           string    `bson:"_id,omitempty"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	FirstName    string    `bson:"first_name"`
	LastName     string    `bson:"last_name"`
	Role         Role      `bson:"role"`
	Phone        string    `bson:"phone,omitempty"`
	Active       bool      `bson:"active"`
	LastLogin    time.Time `bson:"last_login,omitempty"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
	
	// User-centric resource management
	ManagedStores    []UserStore    `bson:"managed_stores,omitempty"`    // Stores this user can manage
	ManagedSuppliers []UserSupplier `bson:"managed_suppliers,omitempty"` // Suppliers this user can manage
}

// UserStore represents a store that a user can manage
type UserStore struct {
	StoreID     string    `bson:"store_id"`
	StoreName   string    `bson:"store_name"`   // Cached for performance
	AccessLevel string    `bson:"access_level"` // READ, WRITE, ADMIN
	AssignedAt  time.Time `bson:"assigned_at"`
	AssignedBy  string    `bson:"assigned_by"`  // User ID who granted access
}

// UserSupplier represents a supplier that a user can manage
type UserSupplier struct {
	SupplierID  string    `bson:"supplier_id"`
	SupplierName string   `bson:"supplier_name"` // Cached for performance
	AccessLevel string    `bson:"access_level"`  // READ, WRITE, ADMIN
	AssignedAt  time.Time `bson:"assigned_at"`
	AssignedBy  string    `bson:"assigned_by"`   // User ID who granted access
}

// NewUser creates a new user
func NewUser(email, password, firstName, lastName string, role Role) (*User, error) {
	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(passwordHash),
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword verifies a password against the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(password string) error {
	// Hash the new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName, phone string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.Phone = phone
	u.UpdatedAt = time.Now()
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now()
}

// RecordLogin records a login event
func (u *User) RecordLogin() {
	u.LastLogin = time.Now()
	u.UpdatedAt = time.Now()
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsStaff checks if the user has staff role
func (u *User) IsStaff() bool {
	return u.Role == RoleStaff || u.Role == RoleAdmin
}

// Address represents a user address
type Address struct {
	ID          string    `bson:"_id,omitempty"`
	UserID      string    `bson:"user_id"`
	Name        string    `bson:"name"`
	Street      string    `bson:"street"`
	City        string    `bson:"city"`
	State       string    `bson:"state"`
	PostalCode  string    `bson:"postal_code"`
	Country     string    `bson:"country"`
	IsDefault   bool      `bson:"is_default"`
	Phone       string    `bson:"phone,omitempty"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// NewAddress creates a new address
func NewAddress(userID, name, street, city, state, postalCode, country, phone string, isDefault bool) *Address {
	now := time.Now()
	return &Address{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        name,
		Street:      street,
		City:        city,
		State:       state,
		PostalCode:  postalCode,
		Country:     country,
		IsDefault:   isDefault,
		Phone:       phone,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the address
func (a *Address) Update(name, street, city, state, postalCode, country, phone string) {
	a.Name = name
	a.Street = street
	a.City = city
	a.State = state
	a.PostalCode = postalCode
	a.Country = country
	a.Phone = phone
	a.UpdatedAt = time.Now()
}

// SetDefault sets this address as the default
func (a *Address) SetDefault(isDefault bool) {
	a.IsDefault = isDefault
	a.UpdatedAt = time.Now()
}

// Store and Supplier Management Methods

// AddManagedStore adds a store to the user's managed stores list
func (u *User) AddManagedStore(storeID, storeName, accessLevel, assignedBy string) {
	// Check if store already exists
	for i, store := range u.ManagedStores {
		if store.StoreID == storeID {
			// Update existing store access
			u.ManagedStores[i].AccessLevel = accessLevel
			u.ManagedStores[i].AssignedBy = assignedBy
			u.ManagedStores[i].AssignedAt = time.Now()
			u.UpdatedAt = time.Now()
			return
		}
	}
	
	// Add new store
	u.ManagedStores = append(u.ManagedStores, UserStore{
		StoreID:     storeID,
		StoreName:   storeName,
		AccessLevel: accessLevel,
		AssignedAt:  time.Now(),
		AssignedBy:  assignedBy,
	})
	u.UpdatedAt = time.Now()
}

// RemoveManagedStore removes a store from the user's managed stores list
func (u *User) RemoveManagedStore(storeID string) {
	for i, store := range u.ManagedStores {
		if store.StoreID == storeID {
			u.ManagedStores = append(u.ManagedStores[:i], u.ManagedStores[i+1:]...)
			u.UpdatedAt = time.Now()
			return
		}
	}
}

// HasStoreAccess checks if user has access to a specific store with given access level
func (u *User) HasStoreAccess(storeID, accessLevel string) bool {
	// Admin has access to all stores
	if u.Role == RoleAdmin {
		return true
	}
	
	for _, store := range u.ManagedStores {
		if store.StoreID == storeID {
			// Check access level hierarchy: ADMIN > WRITE > READ
			switch accessLevel {
			case "READ":
				return store.AccessLevel == "READ" || store.AccessLevel == "WRITE" || store.AccessLevel == "ADMIN"
			case "WRITE":
				return store.AccessLevel == "WRITE" || store.AccessLevel == "ADMIN"
			case "ADMIN":
				return store.AccessLevel == "ADMIN"
			}
		}
	}
	return false
}

// GetManagedStoreIDs returns a list of store IDs the user can manage
func (u *User) GetManagedStoreIDs() []string {
	storeIDs := make([]string, len(u.ManagedStores))
	for i, store := range u.ManagedStores {
		storeIDs[i] = store.StoreID
	}
	return storeIDs
}

// AddManagedSupplier adds a supplier to the user's managed suppliers list
func (u *User) AddManagedSupplier(supplierID, supplierName, accessLevel, assignedBy string) {
	// Check if supplier already exists
	for i, supplier := range u.ManagedSuppliers {
		if supplier.SupplierID == supplierID {
			// Update existing supplier access
			u.ManagedSuppliers[i].AccessLevel = accessLevel
			u.ManagedSuppliers[i].AssignedBy = assignedBy
			u.ManagedSuppliers[i].AssignedAt = time.Now()
			u.UpdatedAt = time.Now()
			return
		}
	}
	
	// Add new supplier
	u.ManagedSuppliers = append(u.ManagedSuppliers, UserSupplier{
		SupplierID:   supplierID,
		SupplierName: supplierName,
		AccessLevel:  accessLevel,
		AssignedAt:   time.Now(),
		AssignedBy:   assignedBy,
	})
	u.UpdatedAt = time.Now()
}

// RemoveManagedSupplier removes a supplier from the user's managed suppliers list
func (u *User) RemoveManagedSupplier(supplierID string) {
	for i, supplier := range u.ManagedSuppliers {
		if supplier.SupplierID == supplierID {
			u.ManagedSuppliers = append(u.ManagedSuppliers[:i], u.ManagedSuppliers[i+1:]...)
			u.UpdatedAt = time.Now()
			return
		}
	}
}

// HasSupplierAccess checks if user has access to a specific supplier with given access level
func (u *User) HasSupplierAccess(supplierID, accessLevel string) bool {
	// Admin has access to all suppliers
	if u.Role == RoleAdmin {
		return true
	}
	
	for _, supplier := range u.ManagedSuppliers {
		if supplier.SupplierID == supplierID {
			// Check access level hierarchy: ADMIN > WRITE > READ
			switch accessLevel {
			case "READ":
				return supplier.AccessLevel == "READ" || supplier.AccessLevel == "WRITE" || supplier.AccessLevel == "ADMIN"
			case "WRITE":
				return supplier.AccessLevel == "WRITE" || supplier.AccessLevel == "ADMIN"
			case "ADMIN":
				return supplier.AccessLevel == "ADMIN"
			}
		}
	}
	return false
}

// GetManagedSupplierIDs returns a list of supplier IDs the user can manage
func (u *User) GetManagedSupplierIDs() []string {
	supplierIDs := make([]string, len(u.ManagedSuppliers))
	for i, supplier := range u.ManagedSuppliers {
		supplierIDs[i] = supplier.SupplierID
	}
	return supplierIDs
}

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

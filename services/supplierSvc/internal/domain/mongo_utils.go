package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDFromString converts a string to a MongoDB ObjectID
func ObjectIDFromString(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// NewObjectID generates a new MongoDB ObjectID
func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

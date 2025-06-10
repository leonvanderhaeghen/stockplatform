package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StringToObjectID converts a string ID to a MongoDB ObjectID
func StringToObjectID(id string) (primitive.ObjectID, error) {
	if id == "" {
		return primitive.NilObjectID, ErrInvalidID
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, ErrInvalidID
	}

	return objectID, nil
}

// IsValidObjectID checks if a string is a valid MongoDB ObjectID
func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

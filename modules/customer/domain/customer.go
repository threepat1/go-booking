package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName         string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName          string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Age               int                `json:"age,omitempty" bson:"age,omitempty"`
	Email             string             `json:"email,omitempty" bson:"email,omitempty"`
	IsVerified        bool               `json:"is_verified,omitempty" bson:"is_verified,omitempty"`
	VerificationToken string             `json:"verification_token,omitempty" bson:"verification_token,omitempty"`
}

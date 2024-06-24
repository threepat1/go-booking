package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Age       int                `json:"age,omitempty" bson:"age,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

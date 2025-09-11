package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Example struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Data      string             `bson:"data" json:"data" binding:"required"`
	CreatedAt *time.Time         `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" binding:"required"`
	CreatedAt *time.Time         `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	User_id   string             `bson:"user_id" json:"user_id"`
}

// CollectionName returns the MongoDB collection name for this model
func (e *Example) CollectionName() string {
	return "examples"
}

func (u *User) CollectionName() string {
	return "users"
}

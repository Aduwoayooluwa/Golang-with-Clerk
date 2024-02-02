package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	FullName string             `bson:"fullName,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Username string             `bson:"username,omitempty"`
}

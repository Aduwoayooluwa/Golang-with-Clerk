package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateUserResponse struct {
	Message string             `json:"message"`
	ID      primitive.ObjectID `json:"_id,omitempty"`
	Email   string             `json:"email"`
}

type BasicUser struct {
	FullName string `bson:"fullName"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
}

type FullUser struct {
	FullName   string `bson:"fullName"`
	Email      string `bson:"email"`
	Username   string `bson:"username"`
	Bio        string `bson:"bio"`
	ProfileUrl string `bson:"profileUrl"`
}

type UpdateUser struct {
	ProfileUrl string `bson:"profileUrl"`
	Bio        string `bson:"bio"`
	UserID     string `json:"userId"`
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Client *mongo.Client
}

func (app *App) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params User

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	collection := app.Client.Database("dailyDB").Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	oid, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert InsertedID to ObjectID")
		return
	}

	respondWithJSON(w, http.StatusAccepted, CreateUserResponse{
		Message: "User created successfully",
		ID:      oid,
		Email:   params.Email,
	})
}

func (app *App) handlerGetUserDetails(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	options := r.URL.Query().Get("options")

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid user ID format %v", err))
		return
	}

	var user User

	err = app.Client.Database("dailyDB").Collection("users").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(w, http.StatusNotFound, "User not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get user details %v", err))
		}
		return
	}

	if options == "basic" {
		respondWithJSON(w, http.StatusOK, BasicUser{
			FullName: user.FullName,
			Email:    user.Email,
			Username: user.Username,
		})
	} else {
		respondWithJSON(w, http.StatusOK, FullUser{
			FullName: user.FullName,
			Email:    user.Email,
			Username: user.Username,
		})
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *App) handleUserUpdateDetails(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var updateParams UpdateUser

	// updating the document format
	update := bson.M{
		"$set": bson.M{
			"bio":             updateParams.Bio,
			"profileImageUrl": updateParams.ProfileUrl,
		},
	}

	if err := decoder.Decode(&updateParams); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	objID, err := primitive.ObjectIDFromHex(updateParams.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	collection := app.Client.Database("dailyDB").Collection("users")

	insertUpdate, err := collection.UpdateByID(context.TODO(), bson.M{"_id": objID}, update)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	if insertUpdate.MatchedCount == 0 {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	respondWithJSON(w, http.StatusOK, AppResponse{
		Message: "User updated successfully",
	})

}

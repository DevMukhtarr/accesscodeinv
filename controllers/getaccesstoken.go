package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/devmukhtarr/accesscodeinv/database"
	"github.com/devmukhtarr/accesscodeinv/middlewares"
	"github.com/devmukhtarr/accesscodeinv/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const prefix = "CHECKCREDIT"

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	randomString, err := generateRandomString(32)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("accesstokens")

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()

	accesscode := prefix + randomString

	newToken := models.AccessCode{
		ID:         primitive.NewObjectID(),
		UserID:     objectID,
		AccessCode: accesscode,
	}

	_, err = collection.InsertOne(ctx, newToken)

	if err != nil {
		http.Error(w, "Failed to create new token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newToken)
}

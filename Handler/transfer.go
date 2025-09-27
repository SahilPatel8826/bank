package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"project/model"
	"project/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertTransferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var transfer model.Transfer
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call repository with transaction
	err := repository.InsertOneTransfer(repository.MongoClient, transfer)
	if err != nil {
		http.Error(w, "Transfer failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func UpdateTransferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.Transfer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := repository.UpdateOneTransfer(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func FetchMyTransfersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get accountID directly from JWT
	accountIDStr, ok := r.Context().Value("accountID").(string)
	if !ok || accountIDStr == "" {
		http.Error(w, "Unauthorized: accountID missing in token", http.StatusUnauthorized)
		return
	}

	// Convert accountID string → ObjectID
	accountID, err := primitive.ObjectIDFromHex(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid accountID format", http.StatusBadRequest)
		return
	}

	// Debug logs (remove in production)
	log.Printf("Fetching transfers for accountID: %s", accountID.Hex())

	transfers, err := repository.GetAllTransferByID(accountID)
	if err != nil {
		log.Printf("Error fetching transfers: %v", err)
		http.Error(w, "Error fetching transfers", http.StatusInternalServerError)
		return
	}

	log.Printf("Transfers found: %d", len(transfers))

	json.NewEncoder(w).Encode(transfers)
}

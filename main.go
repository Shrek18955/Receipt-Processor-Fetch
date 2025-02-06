package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReceiptResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

var (
	receipts = make(map[string]Receipt)
	points   = make(map[string]int)
	mutex    = &sync.Mutex{}
)

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, `{"description": "The receipt is invalid. Please verify input."}`, http.StatusBadRequest)
		return
	}

	if !validateReceipt(receipt) {
		http.Error(w, `{"description": "The receipt is invalid. Please verify input."}`, http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	pts := calculatePoints(receipt)

	mutex.Lock()
	receipts[id] = receipt
	points[id] = pts
	mutex.Unlock()

	response := ReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mutex.Lock()
	pts, exists := points[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, `{"description": "No receipt found for that ID."}`, http.StatusNotFound)
		return
	}

	response := PointsResponse{Points: pts}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

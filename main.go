package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

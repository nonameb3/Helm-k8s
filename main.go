package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"runtime"
	"math"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type LoadResponse struct {
	Status   string `json:"status"`
	Duration int    `json:"duration_ms"`
	CPUCores int    `json:"cpu_cores_used"`
	Message  string `json:"message"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	// Extract duration from URL path /load/:duration
	path := strings.TrimPrefix(r.URL.Path, "/load/")
	duration, err := strconv.Atoi(path)
	if err != nil || duration <= 0 || duration > 60000 {
		http.Error(w, "Invalid duration. Use /load/:milliseconds (1-60000)", http.StatusBadRequest)
		return
	}

	log.Printf("Load test started for %d milliseconds", duration)

	cpuCores := runtime.NumCPU()
	startTime := time.Now()

	// Create CPU-intensive work across all cores
	done := make(chan bool, cpuCores)

	for i := 0; i < cpuCores; i++ {
		go func() {
			defer func() { done <- true }()

			endTime := time.Now().Add(time.Duration(duration) * time.Millisecond)
			for time.Now().Before(endTime) {
				// CPU-intensive calculation
				for j := 0; j < 1000000; j++ {
					math.Sqrt(float64(j))
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < cpuCores; i++ {
		<-done
	}

	actualDuration := int(time.Since(startTime).Milliseconds())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := LoadResponse{
		Status:   "completed",
		Duration: actualDuration,
		CPUCores: cpuCores,
		Message:  "CPU load test completed successfully",
	}

	log.Printf("Load test completed in %d ms using %d CPU cores", actualDuration, cpuCores)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/load/", loadHandler)

	log.Println("Server starting on :8080")
	log.Println("Endpoints available:")
	log.Println("  GET /health - Health check")
	log.Println("  GET /load/:duration - CPU load test (duration in milliseconds, max 60000)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

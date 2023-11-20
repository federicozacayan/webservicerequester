package main

import (
	"bot/actions"
	"fmt"
	"io"
	"net/http"
)

// CORS middleware
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*") // Allow all methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight CORS request
	if r.Method == "OPTIONS" {
		fmt.Printf("OPTIONS - Preflight CORS request for %s\n", r.URL.Path)
		enableCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Enable CORS for the main request
	enableCORS(w)

	// Get the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	fmt.Printf("body: %v\n", body)
	answer := actions.Exec(body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func main() {
	// Configure the handler for GET requests
	http.HandleFunc("/", handler)

	// Start the server on port 8888
	port := 8888
	fmt.Printf("Server listening on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

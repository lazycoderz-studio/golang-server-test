package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// Setup the HTTP server route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a map to hold the response data
		response := map[string]string{
			"host":   r.Host,
			"status": "success",
		}

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("error encoding response: %v", err)
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error while running server:", err)
	}
}

// main.go
package main

import (
	"log"
	"net/http"
	"os"

	"go-cardano-address-safety/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Check for API key
	key := os.Getenv("BLOCKFROST_KEY")
	if key == "" {
		log.Fatal("Set BLOCKFROST_KEY env var")
	}

	// Setup routes
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Page routes
	r.HandleFunc("/", handlers.HandleIndex).Methods("GET")

	// API routes for HTMX
	r.HandleFunc("/analyze", handlers.HandleAnalyze).Methods("POST")

	log.Println("HTG Stack server starting on :8080")
	log.Println("Visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
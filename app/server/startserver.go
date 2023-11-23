package server

import (
	"log"
	"net/http"
	"os"
	"production-go/data"
	"production-go/presentation"

	"github.com/gorilla/handlers"
)

func StartServer() {
	mux := http.NewServeMux()

	pg, err := data.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	presentation.SetupRoutes(pg)

	// Set up CORS middlware
	corsHandler := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-User-ID"}),
		handlers.AllowCredentials(), // This line sets Access-Control-Allow-Credentials to true
	)

	// Wrap the original mux with the CORS handler.
	corsEnabledMux := corsHandler(mux)

	// Use the CORS-enabled mux in your server.
	Port := os.Getenv("PORT")
	log.Printf("Server starting on :%s", Port)
	if err := http.ListenAndServe(":"+Port, corsEnabledMux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

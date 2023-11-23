package server

import (
	"log"
	"net/http"
	"os"
)

func BasicLogging(mux *http.ServeMux) {
	log.Printf("Server starting on :%s", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

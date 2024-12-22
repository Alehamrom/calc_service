package main

import (
	"log"
	"net/http"

	"github.com/Alehamrom/calc_service/internal/handler"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handler.CalculateHandler)

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

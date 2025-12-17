package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KyleKDang/poker-odds-engine/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	router := api.SetupRouter()

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Poker odds engine starting on %s", addr)
	log.Printf("Using Gin Framework")

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

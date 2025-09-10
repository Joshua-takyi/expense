package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Joshua-takyi/expense/server/internal/connection"
	"github.com/Joshua-takyi/expense/server/internal/models"
	"github.com/Joshua-takyi/expense/server/internal/router"
)

func main() {

	if err := connection.InitDb(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := connection.CloseDb(); err != nil {
			log.Fatalf("failed to close database: %v", err)
		}
	}()

	repo := &models.Repository{DB: connection.Client}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Set up the router with the service layer
	r := router.Router(repo)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("failed to run the server: %v", err)
	}
}

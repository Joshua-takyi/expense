package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joshua/expensetracker/internal/connection"
	"github.com/joshua/expensetracker/internal/models"
	"github.com/joshua/expensetracker/internal/router"
)

func main() {
	db, err := connection.InitPsql()
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}
	defer connection.CloseDB(db)

	repo := &models.Repository{DB: db}
	if err := repo.InitTable(); err != nil {
		log.Fatalf("failed to initialize tables: %v", err)
	}

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

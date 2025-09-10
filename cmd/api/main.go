package main

import (
	"fmt"
	"log"

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

	// Set up the router with the service layer
	r := router.Router(repo)
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("failed to run the server: %v", err)
	}
}

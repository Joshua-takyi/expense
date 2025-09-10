package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (r *Repository) InitTable() error {
	sqlFiles, err := filepath.Glob("../../internal/migrations/*.sql") // Changed path to be relative to cmd/api
	if err != nil {
		return fmt.Errorf("failed to find sql files: %w", err)
	}

	log.Printf("Found SQL files: %v\n", sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("No .sql files found, skipping table initialization.")
		return nil
	}

	for _, file := range sqlFiles {
		log.Printf("Attempting to read file: %s\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read sql file %s: %w", file, err)
		}

		log.Printf("Executing SQL from %s:\n%s\n", file, string(content))
		_, err = r.DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute sql file %s: %w", file, err)
		}
		log.Printf("Successfully executed %s\n", file)
	}

	return nil
}

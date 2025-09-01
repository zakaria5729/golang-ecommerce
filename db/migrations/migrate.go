package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/easy-comerce/backend/db"
	"gorm.io/gorm"
)

// RunMigrations executes all SQL migration files
func RunMigrations() error {
	database := db.GetDB()

	// Get the migrations directory
	migrationsDir := "db/migrations"

	// Read all .sql files in the migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsDir, file.Name())
			log.Printf("Running migration: %s", file.Name())

			if err := executeSQLFile(database, filePath); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}

			log.Printf("Migration completed: %s", file.Name())
		}
	}

	return nil
}

// executeSQLFile reads and executes a SQL file
func executeSQLFile(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Split the SQL file into individual statements
	statements := strings.Split(string(content), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		// Execute each SQL statement
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to execute statement: %s, error: %w", statement, err)
		}
	}

	return nil
}

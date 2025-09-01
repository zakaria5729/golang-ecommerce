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

func RunMigrations() error {
	database := db.GetDB()

	migrationsDir := "db/migrations"
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

func executeSQLFile(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	statements := strings.SplitSeq(string(content), ";")

	for statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to execute statement: %s, error: %w", statement, err)
		}
	}

	return nil
}

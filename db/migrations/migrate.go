package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/logger"
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
			logger.Logger.Info("Running migration", "file", file.Name())

			if err := executeSQLFile(database, filePath); err != nil {
				logger.Logger.Error("Failed to execute migration", "error", err, "file", file.Name())
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}

			logger.Logger.Info("Migration completed", "file", file.Name())
		}
	}

	return nil
}

func executeSQLFile(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	statements := strings.Split(string(content), ";")

	for _, statement := range statements {
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

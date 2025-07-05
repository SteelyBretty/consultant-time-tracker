package database

import (
	"log"

	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
)

func Migrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Project{},
		&models.Allocation{},
		&models.TimeEntry{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func CreateIndexes() error {
	log.Println("Creating database indexes...")

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_allocations_week ON allocations(week_starting)").Error; err != nil {
		return err
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_time_entries_date ON time_entries(date)").Error; err != nil {
		return err
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_projects_client ON projects(client_id)").Error; err != nil {
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}

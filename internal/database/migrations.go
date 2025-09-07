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

	// composite unique indexes (match the model tags)
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_client_user_code ON clients(code, user_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_project_user_code ON projects(code, user_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_timeentry_project_user_date ON time_entries(project_id, user_id, date)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_alloc_project_user_week ON allocations(project_id, user_id, week_starting)").Error; err != nil {
		return err
	}

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

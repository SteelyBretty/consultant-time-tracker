package models

import (
	"time"

	"github.com/google/uuid"
)

type Allocation struct {
	BaseModel
	ProjectID    uuid.UUID `gorm:"not null;uniqueIndex:idx_alloc_project_user_week,priority:1" json:"project_id"`
	UserID       uuid.UUID `gorm:"not null;uniqueIndex:idx_alloc_project_user_week,priority:2" json:"user_id"`
	WeekStarting time.Time `gorm:"not null;uniqueIndex:idx_alloc_project_user_week,priority:3" json:"week_starting"`
	Hours        float64   `gorm:"not null" json:"hours"`
	Notes        string    `json:"notes"`
	Project      Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
}

func (Allocation) TableName() string {
	return "allocations"
}

type AllocationKey struct {
	ProjectID    uuid.UUID
	UserID       uuid.UUID
	WeekStarting time.Time
}

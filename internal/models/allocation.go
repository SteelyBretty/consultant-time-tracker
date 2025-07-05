package models

import (
	"time"

	"github.com/google/uuid"
)

type Allocation struct {
	BaseModel
	ProjectID    uuid.UUID `gorm:"not null" json:"project_id"`
	UserID       uuid.UUID `gorm:"not null" json:"user_id"`
	WeekStarting time.Time `gorm:"not null" json:"week_starting"`
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

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TimeEntry struct {
	BaseModel
	ProjectID   uuid.UUID      `gorm:"not null" json:"project_id"`
	UserID      uuid.UUID      `gorm:"not null" json:"user_id"`
	Date        datatypes.Date `gorm:"not null" json:"date"`
	Hours       float64        `gorm:"not null" json:"hours"`
	Description string         `json:"description"`
	IsBillable  bool           `gorm:"default:true" json:"is_billable"`
	Project     Project        `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
}

func (TimeEntry) TableName() string {
	return "time_entries"
}

type TimeEntryKey struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Date      time.Time
}

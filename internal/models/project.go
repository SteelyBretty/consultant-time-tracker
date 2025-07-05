package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusOnHold    ProjectStatus = "on_hold"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

type Project struct {
	BaseModel
	Name         string          `gorm:"not null" json:"name"`
	Code         string          `gorm:"uniqueIndex;not null" json:"code"`
	Description  string          `json:"description"`
	Status       ProjectStatus   `gorm:"default:'active'" json:"status"`
	BillableRate float64         `gorm:"not null" json:"billable_rate"`
	Currency     string          `gorm:"default:'USD'" json:"currency"`
	StartDate    datatypes.Date  `json:"start_date"`
	EndDate      *datatypes.Date `json:"end_date"`
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	ClientID     uuid.UUID       `gorm:"not null" json:"client_id"`
	UserID       uuid.UUID       `gorm:"not null" json:"user_id"`
	Client       Client          `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	User         User            `gorm:"foreignKey:UserID" json:"-"`
	Allocations  []Allocation    `gorm:"foreignKey:ProjectID" json:"allocations,omitempty"`
	TimeEntries  []TimeEntry     `gorm:"foreignKey:ProjectID" json:"time_entries,omitempty"`
}

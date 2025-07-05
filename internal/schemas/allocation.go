package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreateAllocationRequest struct {
	ProjectID    uuid.UUID `json:"project_id" binding:"required"`
	WeekStarting string    `json:"week_starting" binding:"required,datetime=2006-01-02"`
	Hours        float64   `json:"hours" binding:"required,min=0,max=168"`
	Notes        string    `json:"notes" binding:"max=500"`
}

type UpdateAllocationRequest struct {
	Hours float64 `json:"hours" binding:"required,min=0,max=168"`
	Notes string  `json:"notes" binding:"max=500"`
}

type CopyAllocationsRequest struct {
	FromWeek string `json:"from_week" binding:"required,datetime=2006-01-02"`
	ToWeek   string `json:"to_week" binding:"required,datetime=2006-01-02"`
}

type AllocationResponse struct {
	ID           uuid.UUID       `json:"id"`
	ProjectID    uuid.UUID       `json:"project_id"`
	Project      *ProjectSummary `json:"project,omitempty"`
	WeekStarting string          `json:"week_starting"`
	Hours        float64         `json:"hours"`
	Notes        string          `json:"notes"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ProjectSummary struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Code         string        `json:"code"`
	Client       ClientSummary `json:"client"`
	BillableRate float64       `json:"billable_rate"`
	Currency     string        `json:"currency"`
}

type ClientSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

type AllocationListResponse struct {
	Allocations []AllocationResponse `json:"allocations"`
	Total       int64                `json:"total"`
	Offset      int                  `json:"offset"`
	Limit       int                  `json:"limit"`
}

type WeekAllocationResponse struct {
	WeekStarting string               `json:"week_starting"`
	TotalHours   float64              `json:"total_hours"`
	Allocations  []AllocationResponse `json:"allocations"`
}

type CopyAllocationsResponse struct {
	CopiedCount int                  `json:"copied_count"`
	Allocations []AllocationResponse `json:"allocations"`
}

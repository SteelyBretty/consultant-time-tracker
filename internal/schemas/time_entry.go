package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreateTimeEntryRequest struct {
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`
	Date        string    `json:"date" binding:"required,datetime=2006-01-02"`
	Hours       float64   `json:"hours" binding:"required,min=0,max=24"`
	Description string    `json:"description" binding:"required,min=1,max=1000"`
	IsBillable  bool      `json:"is_billable"`
}

type UpdateTimeEntryRequest struct {
	Hours       float64 `json:"hours" binding:"required,min=0,max=24"`
	Description string  `json:"description" binding:"required,min=1,max=1000"`
	IsBillable  bool    `json:"is_billable"`
}

type TimeEntryResponse struct {
	ID          uuid.UUID       `json:"id"`
	ProjectID   uuid.UUID       `json:"project_id"`
	Project     *ProjectSummary `json:"project,omitempty"`
	Date        string          `json:"date"`
	Hours       float64         `json:"hours"`
	Description string          `json:"description"`
	IsBillable  bool            `json:"is_billable"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type TimeEntryListResponse struct {
	TimeEntries []TimeEntryResponse `json:"time_entries"`
	Total       int64               `json:"total"`
	Offset      int                 `json:"offset"`
	Limit       int                 `json:"limit"`
}

type DayEntriesResponse struct {
	Date        string              `json:"date"`
	TotalHours  float64             `json:"total_hours"`
	TimeEntries []TimeEntryResponse `json:"time_entries"`
}

type WeekEntriesResponse struct {
	WeekStarting string              `json:"week_starting"`
	TimeEntries  []TimeEntryResponse `json:"time_entries"`
	DailyTotals  map[string]float64  `json:"daily_totals"`
	WeekTotal    float64             `json:"week_total"`
}

type ProjectWeekComparisonResponse struct {
	ProjectID       uuid.UUID `json:"project_id"`
	WeekStarting    string    `json:"week_starting"`
	AllocatedHours  float64   `json:"allocated_hours"`
	ActualHours     float64   `json:"actual_hours"`
	Variance        float64   `json:"variance"`
	VariancePercent float64   `json:"variance_percent"`
}

type WeekSummaryResponse struct {
	WeekStarting   string               `json:"week_starting"`
	Projects       []ProjectWeekSummary `json:"projects"`
	TotalAllocated float64              `json:"total_allocated"`
	TotalActual    float64              `json:"total_actual"`
}

type ProjectWeekSummary struct {
	ProjectID      uuid.UUID       `json:"project_id"`
	Project        *ProjectSummary `json:"project"`
	AllocatedHours float64         `json:"allocated_hours"`
	ActualHours    float64         `json:"actual_hours"`
	Variance       float64         `json:"variance"`
}

package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name         string    `json:"name" binding:"required,min=1,max=200"`
	Code         string    `json:"code" binding:"required,min=2,max=20,alphanum"`
	Description  string    `json:"description" binding:"max=1000"`
	ClientID     uuid.UUID `json:"client_id" binding:"required"`
	BillableRate float64   `json:"billable_rate" binding:"required,min=0"`
	Currency     string    `json:"currency" binding:"required,len=3"`
	StartDate    string    `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate      *string   `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateProjectRequest struct {
	Name         *string  `json:"name" binding:"omitempty,min=1,max=200"`
	Code         *string  `json:"code" binding:"omitempty,min=2,max=20,alphanum"`
	Description  *string  `json:"description" binding:"omitempty,max=1000"`
	Status       *string  `json:"status" binding:"omitempty,oneof=active on_hold completed cancelled"`
	BillableRate *float64 `json:"billable_rate" binding:"omitempty,min=0"`
	Currency     *string  `json:"currency" binding:"omitempty,len=3"`
	StartDate    *string  `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate      *string  `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	IsActive     *bool    `json:"is_active"`
}

type ProjectResponse struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Code         string          `json:"code"`
	Description  string          `json:"description"`
	Status       string          `json:"status"`
	BillableRate float64         `json:"billable_rate"`
	Currency     string          `json:"currency"`
	StartDate    string          `json:"start_date"`
	EndDate      *string         `json:"end_date,omitempty"`
	IsActive     bool            `json:"is_active"`
	ClientID     uuid.UUID       `json:"client_id"`
	Client       *ClientResponse `json:"client,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int64             `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreateClientRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=200"`
	Code    string `json:"code" binding:"required,min=2,max=20,alphanum"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone" binding:"omitempty,max=50"`
	Address string `json:"address" binding:"omitempty,max=500"`
}

type UpdateClientRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=1,max=200"`
	Code     *string `json:"code" binding:"omitempty,min=2,max=20,alphanum"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone" binding:"omitempty,max=50"`
	Address  *string `json:"address" binding:"omitempty,max=500"`
	IsActive *bool   `json:"is_active"`
}

type ClientResponse struct {
	ID        uuid.UUID               `json:"id"`
	Name      string                  `json:"name"`
	Code      string                  `json:"code"`
	Email     string                  `json:"email"`
	Phone     string                  `json:"phone"`
	Address   string                  `json:"address"`
	IsActive  bool                    `json:"is_active"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Projects  []ClientProjectResponse `json:"projects,omitempty"`
}

type ClientProjectResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Status       string    `json:"status"`
	BillableRate float64   `json:"billable_rate"`
	Currency     string    `json:"currency"`
}

type ClientListResponse struct {
	Clients []ClientResponse `json:"clients"`
	Total   int64            `json:"total"`
	Offset  int              `json:"offset"`
	Limit   int              `json:"limit"`
}

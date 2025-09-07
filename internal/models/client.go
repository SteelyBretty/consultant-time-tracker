package models

import "github.com/google/uuid"

type Client struct {
	BaseModel
	Name     string    `gorm:"not null" json:"name"`
	Code     string    `gorm:"uniqueIndex:idx_client_user_code,priority:1;not null" json:"code"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Address  string    `json:"address"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
	UserID   uuid.UUID `gorm:"not null;uniqueIndex:idx_client_user_code,priority:2" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Projects []Project `gorm:"foreignKey:ClientID" json:"projects,omitempty"`
}

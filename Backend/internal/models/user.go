package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // employee, manager
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Role   string `json:"role"`
	UserID uint   `json:"user_id"`
}

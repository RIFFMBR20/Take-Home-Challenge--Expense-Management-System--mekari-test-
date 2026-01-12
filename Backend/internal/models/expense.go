package models

import "time"

type ExpenseStatus string

const (
	StatusPending      ExpenseStatus = "pending"
	StatusApproved     ExpenseStatus = "approved"
	StatusRejected     ExpenseStatus = "rejected"
	StatusAutoApproved ExpenseStatus = "auto_approved"
	StatusCompleted    ExpenseStatus = "completed"
)

type Expense struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	UserID           uint          `gorm:"not null" json:"user_id"`
	Amount           int64         `gorm:"not null" json:"amount"`
	Description      string        `gorm:"not null" json:"description"`
	ReceiptURL       string        `json:"receipt_url"`
	Status           ExpenseStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ExternalID       string        `gorm:"uniqueIndex" json:"external_id"`
	RequiresApproval bool          `json:"requires_approval"`
	AutoApproved     bool          `json:"auto_approved"`
	SubmittedAt      time.Time     `gorm:"autoCreateTime" json:"submitted_at"`
	ProcessedAt      *time.Time    `json:"processed_at,omitempty"`
	Approvals        []Approval    `gorm:"foreignKey:ExpenseID" json:"approvals,omitempty"`
}

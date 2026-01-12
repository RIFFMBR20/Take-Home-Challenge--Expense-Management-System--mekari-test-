package models

import (
	"time"
)

type Approval struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	ExpenseID  uint          `gorm:"not null" json:"expense_id"`
	ApproverID uint          `gorm:"not null" json:"approver_id"`
	Approver   User          `gorm:"foreignKey:ApproverID" json:"approver"`
	Status     ExpenseStatus `gorm:"type:varchar(20)" json:"status"`
	Notes      string        `json:"notes"`
	CreatedAt  time.Time     `gorm:"autoCreateTime" json:"created_at"`
}

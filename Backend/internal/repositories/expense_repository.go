package repositories

import (
	model "backend-test-mekari/internal/models"

	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *model.Expense) error
	FindByID(id uint) (*model.Expense, error)
	FindAll(userID uint, role string, status string, limit int, offset int) ([]model.Expense, int64, error)
	Update(expense *model.Expense) error
	CreateApproval(approval *model.Approval) error
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db}
}

func (r *expenseRepository) Create(exp *model.Expense) error { return r.db.Create(exp).Error }

func (r *expenseRepository) Update(exp *model.Expense) error {
	return r.db.Model(exp).Select("Status", "ProcessedAt").Updates(exp).Error
}

func (r *expenseRepository) CreateApproval(app *model.Approval) error {
	return r.db.Create(app).Error
}

func (r *expenseRepository) FindByID(id uint) (*model.Expense, error) {
	var exp model.Expense
	err := r.db.Preload("Approvals").First(&exp, id).Error
	return &exp, err
}

func (r *expenseRepository) FindAll(userID uint, role string, status string, limit int, offset int) ([]model.Expense, int64, error) {
	var expenses []model.Expense
	var total int64

	query := r.db.Model(&model.Expense{})

	if role != "manager" {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Limit(limit).Offset(offset).Order("submitted_at desc").Find(&expenses).Error
	return expenses, total, err
}

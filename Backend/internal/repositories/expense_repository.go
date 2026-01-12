package repositories

import (
	models2 "backend-test-mekari/internal/models"

	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *models2.Expense) error
	FindByID(id uint) (*models2.Expense, error)
	FindAll(userID uint, role string, status string, limit int, offset int) ([]models2.Expense, int64, error)
	Update(expense *models2.Expense) error
	CreateApproval(approval *models2.Approval) error
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db}
}

func (r *expenseRepository) Create(exp *models2.Expense) error { return r.db.Create(exp).Error }

func (r *expenseRepository) Update(exp *models2.Expense) error { return r.db.Save(exp).Error }

func (r *expenseRepository) CreateApproval(app *models2.Approval) error {
	return r.db.Create(app).Error
}

func (r *expenseRepository) FindByID(id uint) (*models2.Expense, error) {
	var exp models2.Expense
	err := r.db.Preload("Approvals").First(&exp, id).Error
	return &exp, err
}

func (r *expenseRepository) FindAll(userID uint, role string, status string, limit int, offset int) ([]models2.Expense, int64, error) {
	var expenses []models2.Expense
	var total int64
	query := r.db.Model(&models2.Expense{})

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

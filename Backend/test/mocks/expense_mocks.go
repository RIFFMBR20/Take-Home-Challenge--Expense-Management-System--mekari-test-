package mocks

import (
	"backend-test-mekari/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(exp *models.Expense) error {
	args := m.Called(exp)
	return args.Error(0)
}

func (m *MockExpenseRepository) Update(exp *models.Expense) error {
	args := m.Called(exp)
	return args.Error(0)
}

func (m *MockExpenseRepository) FindByID(id uint) (*models.Expense, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockExpenseRepository) FindAll(userID uint, role string, status string, limit int, offset int) ([]models.Expense, int64, error) {
	args := m.Called(userID, role, status, limit, offset)
	return args.Get(0).([]models.Expense), args.Get(1).(int64), args.Error(2)
}

func (m *MockExpenseRepository) CreateApproval(approval *models.Approval) error {
	args := m.Called(approval)
	return args.Error(0)
}

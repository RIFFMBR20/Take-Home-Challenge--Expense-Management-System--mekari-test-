package test

import (
	"backend-test-mekari/internal/models"
	"backend-test-mekari/internal/services"
	"backend-test-mekari/test/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExpenseService_Submit(t *testing.T) {
	repo := new(mocks.MockExpenseRepository)
	service := services.NewExpenseService(repo)

	t.Run("Success Auto Approved (Amount < 1jt)", func(t *testing.T) {
		expense := &models.Expense{
			Amount:      500000,
			Description: "Break Fast",
		}

		repo.On("Create", mock.Anything).Return(nil).Once()

		result, err := service.Submit(expense)

		assert.NoError(t, err)
		assert.Equal(t, models.StatusAutoApproved, result.Status)
		assert.True(t, result.AutoApproved)
		repo.AssertExpectations(t)
	})

	t.Run("Success Pending (Amount >= 1jt)", func(t *testing.T) {
		expense := &models.Expense{
			Amount:      2000000,
			Description: "Buy Monitor",
		}

		repo.On("Create", mock.Anything).Return(nil).Once()

		result, err := service.Submit(expense)

		assert.NoError(t, err)
		assert.Equal(t, models.StatusPending, result.Status)
		assert.False(t, result.AutoApproved)
		assert.True(t, result.RequiresApproval)
		repo.AssertExpectations(t)
	})

	t.Run("Error Amount Too Low", func(t *testing.T) {
		expense := &models.Expense{
			Amount:      5000,
			Description: "Parking",
		}

		result, err := service.Submit(expense)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "minimum expense amount")
	})

	t.Run("Error Empty Description", func(t *testing.T) {
		expense := &models.Expense{
			Amount: 50000,
		}

		result, err := service.Submit(expense)

		assert.Nil(t, result)
		assert.EqualError(t, err, "description is required")
	})
}

func TestExpenseService_Approve(t *testing.T) {
	repo := new(mocks.MockExpenseRepository)
	service := services.NewExpenseService(repo)

	t.Run("Success Approve", func(t *testing.T) {
		existingExpense := &models.Expense{
			ID:     1,
			Status: models.StatusPending,
			Amount: 2000000,
		}

		repo.On("FindByID", uint(1)).Return(existingExpense, nil).Once()
		repo.On("Update", mock.MatchedBy(func(e *models.Expense) bool {
			return e.Status == models.StatusApproved
		})).Return(nil).Once()
		repo.On("CreateApproval", mock.Anything).Return(nil).Once()

		err := service.Approve(1, 2, "Approved by manager")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed - Not Pending", func(t *testing.T) {
		existingExpense := &models.Expense{
			ID:     1,
			Status: models.StatusRejected,
		}

		repo.On("FindByID", uint(1)).Return(existingExpense, nil).Once()

		err := service.Approve(1, 2, "Notes")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hanya pengajuan pending yang bisa disetujui")
	})

	t.Run("Failed - Not Found", func(t *testing.T) {
		repo.On("FindByID", uint(99)).Return(nil, errors.New("not found")).Once()

		err := service.Approve(99, 2, "Notes")

		assert.Error(t, err)
	})
}

func TestExpenseService_Reject(t *testing.T) {
	repo := new(mocks.MockExpenseRepository)
	service := services.NewExpenseService(repo)

	t.Run("Success Reject", func(t *testing.T) {
		existingExpense := &models.Expense{
			ID: 1,
		}

		repo.On("FindByID", uint(1)).Return(existingExpense, nil).Once()
		repo.On("Update", mock.MatchedBy(func(e *models.Expense) bool {
			return e.Status == models.StatusRejected
		})).Return(nil).Once()
		repo.On("CreateApproval", mock.Anything).Return(nil).Once()

		err := service.Reject(1, 2, "Insufficient documents")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

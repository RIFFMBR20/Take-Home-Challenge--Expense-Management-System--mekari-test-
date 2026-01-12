package services

import (
	model "backend-test-mekari/internal/models"
	"backend-test-mekari/internal/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ExpenseService interface {
	Submit(exp *model.Expense) (*model.Expense, error)
	Approve(id uint, approverID uint, notes string) error
	Reject(id uint, approverID uint, notes string) error
	GetAll(userID uint, role string, status string, limit int, offset int) ([]model.Expense, int64, error)
	GetByID(id uint) (*model.Expense, error)
}

type expenseService struct {
	repo repositories.ExpenseRepository
}

func NewExpenseService(repo repositories.ExpenseRepository) ExpenseService {
	return &expenseService{repo}
}

func (s *expenseService) GetAll(userID uint, role string, status string, limit int, offset int) ([]model.Expense, int64, error) {
	return s.repo.FindAll(userID, role, status, limit, offset)
}

func (s *expenseService) Submit(exp *model.Expense) (*model.Expense, error) {
	if exp.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	if exp.Amount < 10000 {
		return nil, fmt.Errorf("minimum expense amount is IDR 10.000")
	}
	if exp.Amount > 50000000 {
		return nil, fmt.Errorf("maximum expense amount is IDR 50.000.000")
	}

	exp.ExternalID = uuid.New().String()

	if exp.Amount < 1000000 {
		exp.Status = model.StatusAutoApproved
		exp.AutoApproved = true
		exp.RequiresApproval = false
	} else {
		exp.Status = model.StatusPending
		exp.AutoApproved = false
		exp.RequiresApproval = true
	}

	if err := s.repo.Create(exp); err != nil {
		return nil, err
	}

	if exp.AutoApproved {
		go s.ProcessPayment(exp)
	}

	return exp, nil
}

func (s *expenseService) Approve(id uint, approverID uint, notes string) error {
	exp, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if exp.Status != model.StatusPending {
		return fmt.Errorf("hanya pengajuan pending yang bisa disetujui")
	}

	exp.Status = model.StatusApproved
	if err := s.repo.Update(exp); err != nil {
		return err
	}

	err = s.repo.CreateApproval(&model.Approval{
		ExpenseID:  id,
		ApproverID: approverID,
		Status:     model.StatusApproved,
		Notes:      notes,
	})
	if err != nil {
		return err
	}

	go s.ProcessPayment(exp)
	return nil
}

func (s *expenseService) Reject(id uint, approverID uint, notes string) error {
	exp, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	exp.Status = model.StatusRejected
	if err := s.repo.Update(exp); err != nil {
		return err
	}

	return s.repo.CreateApproval(&model.Approval{
		ExpenseID:  id,
		ApproverID: approverID,
		Status:     model.StatusRejected,
		Notes:      notes,
	})
}

func (s *expenseService) ProcessPayment(exp *model.Expense) {
	url := "https://1620e98f-7759-431c-a2aa-f449d591150b.mock.pstmn.io/v1/payments"

	payload := map[string]interface{}{
		"amount":      exp.Amount,
		"external_id": exp.ExternalID,
	}

	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Printf("Payment failed for ID %d: %v\n", exp.ID, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {
		var result map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return
		}

		if data, ok := result["data"].(map[string]interface{}); ok {
			if data["status"] == "success" {
				exp.Status = model.StatusCompleted
				now := time.Now()
				exp.ProcessedAt = &now
				err := s.repo.Update(exp)
				if err != nil {
					return
				}
			}
		}
	}
}

func (s *expenseService) GetByID(id uint) (*model.Expense, error) {
	exp, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("expense not found")
	}
	return exp, nil
}

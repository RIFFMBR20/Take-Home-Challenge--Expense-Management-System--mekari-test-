package controllers

import (
	"backend-test-mekari/internal/middleware"
	"backend-test-mekari/internal/models"
	"backend-test-mekari/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ExpenseController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Approve(w http.ResponseWriter, r *http.Request)
	HandleList(w http.ResponseWriter, r *http.Request)
	Reject(w http.ResponseWriter, r *http.Request)
	GetDetail(w http.ResponseWriter, r *http.Request)
}

type expenseController struct {
	svc services.ExpenseService
}

func NewExpenseController(svc services.ExpenseService) ExpenseController {
	return &expenseController{svc}
}

// Create godoc
// @Summary Submit a new expense
// @Tags expenses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param expense body models.Expense true "Expense object"
// @Success 201 {object} models.Expense
// @Router /api/expenses [post]
func (c *expenseController) Create(w http.ResponseWriter, r *http.Request) {
	var exp models.Expense
	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		c.renderError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	val := r.Context().Value(middleware.UserKey)
	if user, ok := val.(middleware.AuthUser); ok {
		exp.UserID = user.ID
	} else {
		c.renderError(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}

	res, err := c.svc.Submit(&exp)
	if err != nil {
		c.renderError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// HandleList godoc
// @Summary Get list of expenses
// @Tags expenses
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter status"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} map[string]interface{}
// @Router /api/expenses [get]
func (c *expenseController) HandleList(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(middleware.UserKey)
	user, ok := val.(middleware.AuthUser)
	if !ok {
		c.renderError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	expenses, total, err := c.svc.GetAll(user.ID, user.Role, status, limit, offset)
	if err != nil {
		c.renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": expenses,
		"meta": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetDetail godoc
// @Summary Get expense detail
// @Tags expenses
// @Security BearerAuth
// @Param id path int true "Expense ID"
// @Success 200 {object} models.Expense
// @Router /api/expenses/{id} [get]
func (c *expenseController) GetDetail(w http.ResponseWriter, r *http.Request) {
	id := c.parseID(r.URL.Path)
	if id == 0 {
		c.renderError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := c.svc.GetByID(id)
	if err != nil {
		c.renderError(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Approve godoc
// @Summary Approve an expense
// @Tags expenses
// @Security BearerAuth
// @Param id path int true "Expense ID"
// @Success 204 "No Content"
// @Router /api/expenses/{id}/approve [put]
func (c *expenseController) Approve(w http.ResponseWriter, r *http.Request) {
	id := c.parseID(r.URL.Path)

	val := r.Context().Value(middleware.UserKey)
	user := val.(middleware.AuthUser)

	if err := c.svc.Approve(id, user.ID, "Approved via API"); err != nil {
		c.renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Reject godoc
// @Summary Reject an expense
// @Tags expenses
// @Security BearerAuth
// @Param id path int true "Expense ID"
// @Success 204 "No Content"
// @Router /api/expenses/{id}/reject [put]
func (c *expenseController) Reject(w http.ResponseWriter, r *http.Request) {
	id := c.parseID(r.URL.Path)
	if err := c.svc.Reject(id, 2, "Rejected via API"); err != nil {
		c.renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *expenseController) parseID(path string) uint {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0
	}
	id, _ := strconv.Atoi(parts[2])
	return uint(id)
}

func (c *expenseController) renderError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

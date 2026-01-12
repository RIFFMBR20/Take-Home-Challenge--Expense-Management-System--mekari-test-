package routers

import (
	"backend-test-mekari/internal/controllers"
	"backend-test-mekari/internal/middleware"
	"net/http"
	"strings"
)

func SetupRouter(ctrl controllers.ExpenseController, authCtrl *controllers.AuthController) *http.ServeMux {
	mux := http.NewServeMux()

	// 1. POST /api/auth/login
	mux.HandleFunc("/api/auth/login", authCtrl.HandleLogin)

	// 2. GET /api/health
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"UP"}`))
	})

	// 3. GET & POST /api/expenses
	mux.HandleFunc("/api/expenses", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctrl.Create(w, r)
		} else if r.Method == http.MethodGet {
			ctrl.HandleList(w, r)
		}
	}))

	// 4. Endpoint with ID (:id)
	mux.HandleFunc("/api/expenses/", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/expenses/")
		segments := strings.Split(path, "/")

		// GET /api/expenses/:id
		if len(segments) == 1 && r.Method == http.MethodGet {
			ctrl.GetDetail(w, r)
			return
		}

		// PUT /api/expenses/:id/approve
		if len(segments) == 2 && segments[1] == "approve" && r.Method == http.MethodPut {
			ctrl.Approve(w, r)
			return
		}

		// PUT /api/expenses/:id/reject
		if len(segments) == 2 && segments[1] == "reject" && r.Method == http.MethodPut {
			ctrl.Reject(w, r)
			return
		}

		http.NotFound(w, r)
	}))

	return mux
}

package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const UserKey contextKey = "user"

type AuthUser struct {
	ID   uint
	Role string
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.Header.Get("X-User-ID")
		role := r.Header.Get("X-User-Role")

		id, _ := strconv.Atoi(idStr)
		if id == 0 {
			id = 1
		}
		if role == "" {
			role = "employee"
		}

		user := AuthUser{
			ID:   uint(id),
			Role: role,
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

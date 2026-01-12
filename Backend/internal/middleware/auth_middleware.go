package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserKey contextKey = "user"

type AuthUser struct {
	ID   uint
	Role string
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := AuthUser{ID: 1, Role: "manager"}

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

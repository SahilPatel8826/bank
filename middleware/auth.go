package middleware

import (
	"context"
	"net/http"
	"strings"

	"project/utils"
)

func Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// ✅ extract claims safely
		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized: userID missing in token", http.StatusUnauthorized)
			return
		}

		accountID, ok := claims["accountID"].(string)
		if !ok || accountID == "" {
			http.Error(w, "Unauthorized: accountID missing in token", http.StatusUnauthorized)
			return
		}

		role, _ := claims["role"].(string)
		email, _ := claims["email"].(string)

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "accountID", accountID)

		next(w, r.WithContext(ctx))
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

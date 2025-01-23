package middleware

import (
	"DeliFood/backend/pkg/repo"
	"net/http"
	"strings"
)

func RoleMiddleware(userRepo *repo.UserRepo, allowedRoles ...string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := userRepo.GetUserByToken(token)
			if err != nil || user.Token != token {
				http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Check if the user's role is allowed
			for _, role := range allowedRoles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		})
	}
}

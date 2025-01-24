package middleware

import (
	"DeliFood/backend/utils"
	"net/http"
	"strings"
)

func RoleMiddleware(allowedRoles ...string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateJWT(token)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		})
	}
}

package middleware

import (
	"DeliFood/backend/pkg/repo"
	"log"
	"net/http"
	"strings"
)

func RoleMiddleware(userRepo *repo.UserRepo, allowedRoles ...string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the Authorization header
			authHeader := r.Header.Get("Authorization")
			log.Printf("Authorization Header: %s", authHeader)
			if authHeader == "" {
				log.Println("Unauthorized: No Authorization header provided")
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}

			// Extract the token from the header
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				log.Println("Unauthorized: Empty token provided")
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			// Fetch the user using the token
			user, err := userRepo.GetUserByToken(token)
			if err != nil {
				log.Printf("Unauthorized: Invalid or expired token: %v", err)
				http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Validate the token matches the user's stored token
			if user.Token != token {
				log.Printf("Unauthorized: Token mismatch for user %s", user.Email)
				http.Error(w, "Unauthorized: Token mismatch", http.StatusUnauthorized)
				return
			}

			// Check if the user's role is in the list of allowed roles
			for _, role := range allowedRoles {
				if user.Role == role {
					log.Printf("Access granted to user: %s with role: %s", user.Email, user.Role)
					next.ServeHTTP(w, r)
					return
				}
			}

			// If the user's role is not allowed, block access
			log.Printf("Forbidden: User %s with role %s attempted to access a restricted route", user.Email, user.Role)
			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		})
	}
}

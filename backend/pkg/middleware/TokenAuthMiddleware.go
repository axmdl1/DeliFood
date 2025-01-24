package middleware

/*
func TokenAuthMiddleware(userRepo *repo.UserRepo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Extract token from the header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Verify token in the database
		user, err := userRepo.GetUserByToken(token)
		if err != nil || user == nil {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
*/

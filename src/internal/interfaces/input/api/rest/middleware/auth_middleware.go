package middleware

import (
	"go-project/src/internal/adaptors/ports" // import your ports interface
	"net/http"
)

type AuthMiddleware struct {
	UserRepo ports.UserRepository
}

func NewAuthMiddleware(userRepo ports.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{UserRepo: userRepo}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Missing Authentication cookie", http.StatusUnauthorized)
			return
		}

		// Check token in backend (DB)
		ctx := r.Context()
		token := cookie.Value
		valid, err := m.UserRepo.IsTokenValid(ctx, token)
		if err != nil || !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid; pass request forward
		next.ServeHTTP(w, r)
	})
}

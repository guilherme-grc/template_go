package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"template.go/internal/auth"
)

// Auth — equivalent to Laravel's middleware('auth:sanctum') or middleware('auth:api')
// Protects routes that require authentication
func Auth(jwtSvc *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization: Bearer <token> header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "authorization token not provided")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				respondUnauthorized(w, "invalid format. Use: Bearer <token>")
				return
			}

			// Using the translated ValidateToken method
			claims, err := jwtSvc.ValidateToken(parts[1])
			if err != nil {
				respondUnauthorized(w, "invalid or expired token")
				return
			}

			// Ensures it is an access token (not refresh)
			if claims.TokenType != auth.AccessToken {
				respondUnauthorized(w, "please use an access token for authentication")
				return
			}

			// Injects user into context — equivalent to auth()->user()
			ctx := auth.InjectUserIntoContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Guest — equivalent to Laravel's middleware('guest')
// Blocks already authenticated users (e.g., login/register routes)
func Guest(jwtSvc *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 {
					if _, err := jwtSvc.ValidateToken(parts[1]); err == nil {
						respondJSON(w, http.StatusForbidden, map[string]string{
							"message": "you are already authenticated",
						})
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func respondUnauthorized(w http.ResponseWriter, msg string) {
	respondJSON(w, http.StatusUnauthorized, map[string]string{"message": msg})
}

func respondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

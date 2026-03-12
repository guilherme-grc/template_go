package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"reembolso/internal/auth"
)

// Auth — equivalente ao middleware('auth:sanctum') ou middleware('auth:api') do Laravel
// Protege rotas que exigem autenticação
func Auth(jwtSvc *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai o token do header Authorization: Bearer <token>
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "token não fornecido")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				respondUnauthorized(w, "formato inválido. Use: Bearer <token>")
				return
			}

			claims, err := jwtSvc.ValidarToken(parts[1])
			if err != nil {
				respondUnauthorized(w, "token inválido ou expirado")
				return
			}

			// Garante que é um access token (não refresh)
			if claims.TokenType != auth.AccessToken {
				respondUnauthorized(w, "use o access token para autenticar")
				return
			}

			// Injeta o usuário no contexto — equivalente ao auth()->user()
			ctx := auth.InjetarUsuarioNoContexto(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Guest — equivalente ao middleware('guest') do Laravel
// Bloqueia usuários já autenticados (ex: rota de login)
func Guest(jwtSvc *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 {
					if _, err := jwtSvc.ValidarToken(parts[1]); err == nil {
						respondJSON(w, http.StatusForbidden, map[string]string{
							"message": "você já está autenticado",
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

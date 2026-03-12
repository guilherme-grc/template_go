package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins []string // equivalente ao 'allowed_origins' do config/cors.php
	AllowedMethods []string // equivalente ao 'allowed_methods'
	AllowedHeaders []string // equivalente ao 'allowed_headers'
	MaxAge         int      // equivalente ao 'max_age'
}

// DefaultCORSConfig — padrão semelhante ao Laravel
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		MaxAge:         86400,
	}
}

// CORS — equivalente ao middleware HandleCors do Laravel (fruitcake/laravel-cors)
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if isAllowedOrigin(origin, cfg.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if contains(cfg.AllowedOrigins, "*") {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", http.StatusText(cfg.MaxAge))
			w.Header().Set("Vary", "Origin")

			// Preflight request (OPTIONS) — responde imediatamente
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAllowedOrigin(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

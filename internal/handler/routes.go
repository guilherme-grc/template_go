package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"template.go/internal/auth"
	"template.go/internal/middleware"
	"template.go/internal/service"
)

// RegisterRoutes - Orchestrates route definitions and middleware application
func RegisterRoutes(
	mux *http.ServeMux,
	authSvc *service.AuthService,
	jwtSvc *auth.JWTService,
) http.Handler {
	authHandler := NewAuthHandler(authSvc)

	authMiddleware := middleware.Auth(jwtSvc)
	guestMiddleware := middleware.Guest(jwtSvc)

	// throttle:60,1 and throttle:10,1 — equivalent to Laravel's Rate Limiting
	globalLimiter := middleware.NewRateLimiter(60, time.Minute)
	authLimiter := middleware.NewRateLimiter(10, time.Minute)

	mux.HandleFunc("/health", healthCheck)

	// Auth Routes
	mux.Handle("/auth/register", middleware.Chain(
		http.HandlerFunc(authHandler.Register),
		guestMiddleware,
		authLimiter.Throttle,
	))

	mux.Handle("/auth/login", middleware.Chain(
		http.HandlerFunc(authHandler.Login),
		guestMiddleware,
		authLimiter.Throttle,
	))

	mux.Handle("/auth/refresh", http.HandlerFunc(authHandler.Refresh))

	mux.Handle("/auth/me", middleware.Chain(
		http.HandlerFunc(authHandler.Me),
		authMiddleware,
	))

	mux.Handle("/auth/logout", middleware.Chain(
		http.HandlerFunc(authHandler.Logout),
		authMiddleware,
	))

	// Global Middlewares — equivalent to Laravel's bootstrap/app.php (Global Middleware Stack)
	return middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.RequestLogger,
		middleware.CORS(middleware.DefaultCORSConfig()),
		globalLimiter.Throttle,
	)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

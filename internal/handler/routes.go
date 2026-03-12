package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"reembolso/internal/auth"
	"reembolso/internal/middleware"
	"reembolso/internal/service"
)

func RegisterRoutes(
	mux *http.ServeMux,
	authSvc *service.AuthService,
	reembolsoSvc *service.ReembolsoService,
	jwtSvc *auth.JWTService,
) http.Handler {
	authHandler := NewAuthHandler(authSvc)
	reembolsoHandler := NewReembolsoHandler(reembolsoSvc)

	authMiddleware := middleware.Auth(jwtSvc)
	guestMiddleware := middleware.Guest(jwtSvc)

	// throttle:60,1 e throttle:10,1 — igual ao Laravel
	globalLimiter := middleware.NewRateLimiter(60, time.Minute)
	authLimiter := middleware.NewRateLimiter(10, time.Minute)

	mux.HandleFunc("/health", healthCheck)

	mux.Handle("/auth/register", middleware.Chain(
		http.HandlerFunc(authHandler.Register),
		guestMiddleware, authLimiter.Throttle,
	))
	mux.Handle("/auth/login", middleware.Chain(
		http.HandlerFunc(authHandler.Login),
		guestMiddleware, authLimiter.Throttle,
	))
	mux.Handle("/auth/refresh", http.HandlerFunc(authHandler.Refresh))

	mux.Handle("/auth/me", middleware.Chain(
		http.HandlerFunc(authHandler.Me), authMiddleware,
	))
	mux.Handle("/auth/logout", middleware.Chain(
		http.HandlerFunc(authHandler.Logout), authMiddleware,
	))
	mux.Handle("/reembolsos", middleware.Chain(
		http.HandlerFunc(reembolsoHandler.Criar), authMiddleware,
	))
	mux.Handle("/reembolsos/", middleware.Chain(
		http.HandlerFunc(reembolsoHandler.Roteador), authMiddleware,
	))

	// Middlewares globais — equivalente ao bootstrap/app.php do Laravel
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

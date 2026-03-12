package middleware

import "net/http"

// Chain — encadeia múltiplos middlewares
// Equivalente ao Route::middleware(['auth', 'verified']) do Laravel
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

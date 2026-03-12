package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type client struct {
	count    int
	resetAt  time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	limit    int
	window   time.Duration
}

// NewRateLimiter — cria um limiter
// Equivalente ao Route::middleware('throttle:60,1') do Laravel
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
	// Limpeza periódica de clientes expirados
	go rl.cleanup()
	return rl
}

// Throttle — middleware que aplica o rate limiting
// Equivalente ao middleware('throttle:limite,minutos') do Laravel
func (rl *RateLimiter) Throttle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)

		rl.mu.Lock()
		c, exists := rl.clients[ip]
		now := time.Now()

		if !exists || now.After(c.resetAt) {
			c = &client{count: 0, resetAt: now.Add(rl.window)}
			rl.clients[ip] = c
		}

		c.count++
		remaining := rl.limit - c.count
		resetIn := int(time.Until(c.resetAt).Seconds())
		rl.mu.Unlock()

		// Headers padrão do Laravel
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(0, remaining)))
		w.Header().Set("X-RateLimit-Reset", strconv.Itoa(resetIn))

		if remaining < 0 {
			w.Header().Set("Retry-After", strconv.Itoa(resetIn))
			writeError(w, http.StatusTooManyRequests, "Muitas requisições. Tente novamente em breve.", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	for range time.Tick(5 * time.Minute) {
		rl.mu.Lock()
		now := time.Now()
		for ip, c := range rl.clients {
			if now.After(c.resetAt) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func extractIP(r *http.Request) string {
	// Respeita X-Forwarded-For (proxy/load balancer), igual ao Laravel
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	return r.RemoteAddr
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

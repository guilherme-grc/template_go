package middleware

import (
	"encoding/json"
	"net/http"

	"template.go/internal/logger"
	"template.go/internal/validation"
)

// ErrorResponse — formato padrão de erro, igual ao Laravel
type ErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

// Recovery — equivalente ao Handler::render() do Laravel
// Captura panics e transforma em respostas JSON estruturadas
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic capturado", logger.With("error", err), logger.With("path", r.URL.Path))
				writeError(w, http.StatusInternalServerError, "Erro interno do servidor.", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequestLogger — equivalente ao middleware de log HTTP do Laravel
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Info("request", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": rw.status,
			"ip":     r.RemoteAddr,
		})
	})
}

// HandleValidation — helper para tratar ValidationError e responder no formato Laravel
func HandleValidation(w http.ResponseWriter, err error) bool {
	if ve, ok := err.(*validation.ValidationError); ok {
		writeError(w, http.StatusUnprocessableEntity, "Os dados fornecidos são inválidos.", ve.Errors)
		return true
	}
	return false
}

func writeError(w http.ResponseWriter, status int, msg string, errs map[string][]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Message: msg, Errors: errs})
}

// responseWriter captura o status code para logging
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

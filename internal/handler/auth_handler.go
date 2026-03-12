package handler

import (
	"encoding/json"
	"net/http"

	"reembolso/internal/auth"
	"reembolso/internal/middleware"
	"reembolso/internal/model"
	"reembolso/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "body inválido")
		return
	}
	usuario, tokens, err := h.authSvc.Register(req)
	if err != nil {
		if middleware.HandleValidation(w, err) {
			return
		}
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{"usuario": usuario, "tokens": tokens})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "body inválido")
		return
	}
	tokens, err := h.authSvc.Login(req)
	if err != nil {
		if middleware.HandleValidation(w, err) {
			return
		}
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "body inválido")
		return
	}
	tokens, err := h.authSvc.Refresh(req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.UsuarioDoContexto(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "não autenticado")
		return
	}
	usuario, err := h.authSvc.Me(claims.UsuarioID)
	if err != nil {
		respondError(w, http.StatusNotFound, "usuário não encontrado")
		return
	}
	respondJSON(w, http.StatusOK, usuario)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "logout realizado com sucesso. Descarte o token no cliente.",
	})
}

func respondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"message": msg})
}

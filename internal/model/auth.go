package model

// LoginRequest — body do POST /auth/login
type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// RegisterRequest — body do POST /auth/register
type RegisterRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// RefreshRequest — body do POST /auth/refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

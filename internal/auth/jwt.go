package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims — equivalente ao payload do token no Laravel
type Claims struct {
	UsuarioID int64     `json:"usuario_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // segundos
}

type JWTService struct {
	secret            []byte
	accessExpiryMin   int
	refreshExpiryDays int
}

func NewJWTService(secret string, accessMin, refreshDays int) *JWTService {
	return &JWTService{
		secret:            []byte(secret),
		accessExpiryMin:   accessMin,
		refreshExpiryDays: refreshDays,
	}
}

// GerarTokens — equivalente ao Auth::attempt() + token() do Laravel
func (j *JWTService) GerarTokens(usuarioID int64, email string) (*TokenPair, error) {
	accessToken, err := j.gerarToken(usuarioID, email, AccessToken, time.Duration(j.accessExpiryMin)*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.gerarToken(usuarioID, email, RefreshToken, time.Duration(j.refreshExpiryDays)*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    j.accessExpiryMin * 60,
	}, nil
}

func (j *JWTService) gerarToken(usuarioID int64, email string, tokenType TokenType, expiry time.Duration) (string, error) {
	claims := Claims{
		UsuarioID: usuarioID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidarToken — equivalente ao middleware auth do Laravel
func (j *JWTService) ValidarToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// ValidarRefreshToken — valida especificamente um refresh token
func (j *JWTService) ValidarRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := j.ValidarToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != RefreshToken {
		return nil, errors.New("token não é um refresh token")
	}
	return claims, nil
}

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

// Claims - Equivalent to the token payload in Laravel
type Claims struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // in seconds
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

// GenerateTokenPair - Equivalent to Laravel's Auth::attempt() + token()
func (j *JWTService) GenerateTokenPair(userID int64, email string) (*TokenPair, error) {
	accessToken, err := j.generateToken(userID, email, AccessToken, time.Duration(j.accessExpiryMin)*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateToken(userID, email, RefreshToken, time.Duration(j.refreshExpiryDays)*24*time.Hour)
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

func (j *JWTService) generateToken(userID int64, email string, tokenType TokenType, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
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

// ValidateToken - Equivalent to Laravel's 'auth' middleware
func (j *JWTService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken - Specifically validates a refresh token
func (j *JWTService) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != RefreshToken {
		return nil, errors.New("token is not a refresh token")
	}
	return claims, nil
}

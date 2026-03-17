package service

import (
	"errors"

	"template.go/internal/auth"
	"template.go/internal/model"
	"template.go/internal/repository"
	"template.go/internal/validation"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtSvc   *auth.JWTService
}

func NewAuthService(repo *repository.UserRepository, jwtSvc *auth.JWTService) *AuthService {
	return &AuthService{userRepo: repo, jwtSvc: jwtSvc}
}

// Register with validation — equivalent to Laravel's FormRequest/RegisterController
func (s *AuthService) Register(req model.RegisterRequest) (*model.User, *auth.TokenPair, error) {
	v := validation.New(map[string]interface{}{
		"name":     req.Name,
		"email":    req.Email,
		"password": req.Password,
	})

	if err := v.Validate(map[string][]string{
		"name":     {"required", "min:2", "max:100"},
		"email":    {"required", "email", "max:150"},
		"password": {"required", "min:8"},
	}); err != nil {
		return nil, nil, err
	}

	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.New("email already registered")
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.Create(req.Name, req.Email, passwordHash)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.jwtSvc.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Login with validation
func (s *AuthService) Login(req model.LoginRequest) (*auth.TokenPair, error) {
	v := validation.New(map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
	})

	if err := v.Validate(map[string][]string{
		"email":    {"required", "email"},
		"password": {"required"},
	}); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Note: using the translated CheckPassword function
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.jwtSvc.GenerateTokenPair(user.ID, user.Email)
}

func (s *AuthService) Refresh(refreshToken string) (*auth.TokenPair, error) {
	claims, err := s.jwtSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.jwtSvc.GenerateTokenPair(user.ID, user.Email)
}

func (s *AuthService) Me(userID int64) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

package service

import (
	"errors"

	"reembolso/internal/auth"
	"reembolso/internal/model"
	"reembolso/internal/repository"
	"reembolso/internal/validation"
)

type AuthService struct {
	usuarioRepo *repository.UsuarioRepository
	jwtSvc      *auth.JWTService
}

func NewAuthService(repo *repository.UsuarioRepository, jwtSvc *auth.JWTService) *AuthService {
	return &AuthService{usuarioRepo: repo, jwtSvc: jwtSvc}
}

// Register com validação — equivalente ao FormRequest do Laravel
func (s *AuthService) Register(req model.RegisterRequest) (*model.Usuario, *auth.TokenPair, error) {
	v := validation.New(map[string]interface{}{
		"nome":  req.Nome,
		"email": req.Email,
		"senha": req.Senha,
	})
	if err := v.Validate(map[string][]string{
		"nome":  {"required", "min:2", "max:100"},
		"email": {"required", "email", "max:150"},
		"senha": {"required", "min:8"},
	}); err != nil {
		return nil, nil, err
	}

	existe, err := s.usuarioRepo.EmailExiste(req.Email)
	if err != nil {
		return nil, nil, err
	}
	if existe {
		return nil, nil, errors.New("email já cadastrado")
	}

	senhaHash, err := auth.HashSenha(req.Senha)
	if err != nil {
		return nil, nil, err
	}

	usuario, err := s.usuarioRepo.Criar(req.Nome, req.Email, senhaHash)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.jwtSvc.GerarTokens(usuario.ID, usuario.Email)
	if err != nil {
		return nil, nil, err
	}

	return usuario, tokens, nil
}

// Login com validação
func (s *AuthService) Login(req model.LoginRequest) (*auth.TokenPair, error) {
	v := validation.New(map[string]interface{}{
		"email": req.Email,
		"senha": req.Senha,
	})
	if err := v.Validate(map[string][]string{
		"email": {"required", "email"},
		"senha": {"required"},
	}); err != nil {
		return nil, err
	}

	usuario, err := s.usuarioRepo.BuscarPorEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	if !auth.VerificarSenha(req.Senha, usuario.SenhaHash) {
		return nil, errors.New("credenciais inválidas")
	}

	return s.jwtSvc.GerarTokens(usuario.ID, usuario.Email)
}

func (s *AuthService) Refresh(refreshToken string) (*auth.TokenPair, error) {
	claims, err := s.jwtSvc.ValidarRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token inválido ou expirado")
	}

	usuario, err := s.usuarioRepo.BuscarPorID(claims.UsuarioID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	return s.jwtSvc.GerarTokens(usuario.ID, usuario.Email)
}

func (s *AuthService) Me(usuarioID int64) (*model.Usuario, error) {
	return s.usuarioRepo.BuscarPorID(usuarioID)
}

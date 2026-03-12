package service

import (
	"errors"

	"reembolso/internal/model"
	"reembolso/internal/repository"
	"reembolso/internal/validation"
)

type ReembolsoService struct {
	repo *repository.ReembolsoRepository
}

func NewReembolsoService(repo *repository.ReembolsoRepository) *ReembolsoService {
	return &ReembolsoService{repo: repo}
}

// CriarReembolso com validação — equivalente ao FormRequest do Laravel
func (s *ReembolsoService) CriarReembolso(req model.CriarReembolsoRequest) (*model.Reembolso, error) {
	v := validation.New(map[string]interface{}{
		"descricao": req.Descricao,
		"valor":     req.Valor,
		"categoria": req.Categoria,
	})
	if err := v.Validate(map[string][]string{
		"descricao": {"required", "min:3", "max:500"},
		"valor":     {"required", "numeric", "gt:0"},
		"categoria": {"required", "in:TRANSPORTE,HOSPEDAGEM,ALIMENTACAO,MATERIAL,OUTROS"},
	}); err != nil {
		return nil, err
	}

	return s.repo.Criar(req)
}

func (s *ReembolsoService) BuscarReembolso(id int64) (*model.Reembolso, error) {
	return s.repo.BuscarPorID(id)
}

func (s *ReembolsoService) ListarReembolsos(usuarioID int64) ([]model.Reembolso, error) {
	return s.repo.ListarPorUsuario(usuarioID)
}

func (s *ReembolsoService) AprovarReembolso(id int64) error {
	reembolso, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	if reembolso.Status != model.StatusPendente {
		return errors.New("apenas reembolsos pendentes podem ser aprovados")
	}
	return s.repo.AtualizarStatus(id, model.StatusAprovado)
}

func (s *ReembolsoService) RejeitarReembolso(id int64) error {
	reembolso, err := s.repo.BuscarPorID(id)
	if err != nil {
		return err
	}
	if reembolso.Status != model.StatusPendente {
		return errors.New("apenas reembolsos pendentes podem ser rejeitados")
	}
	return s.repo.AtualizarStatus(id, model.StatusRejeitado)
}

package repository

import (
	"database/sql"
	"time"

	"reembolso/internal/model"
)

type ReembolsoRepository struct {
	db *sql.DB
}

func NewReembolsoRepository(db *sql.DB) *ReembolsoRepository {
	return &ReembolsoRepository{db: db}
}

func (r *ReembolsoRepository) Criar(req model.CriarReembolsoRequest) (*model.Reembolso, error) {
	query := `
		INSERT INTO reembolsos (usuario_id, descricao, valor, categoria, status, comprovante, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	agora := time.Now()
	reembolso := &model.Reembolso{
		UsuarioID:    req.UsuarioID,
		Descricao:    req.Descricao,
		Valor:        req.Valor,
		Categoria:    req.Categoria,
		Status:       model.StatusPendente,
		Comprovante:  req.Comprovante,
		CriadoEm:    agora,
		AtualizadoEm: agora,
	}

	err := r.db.QueryRow(query,
		reembolso.UsuarioID,
		reembolso.Descricao,
		reembolso.Valor,
		reembolso.Categoria,
		reembolso.Status,
		reembolso.Comprovante,
		reembolso.CriadoEm,
		reembolso.AtualizadoEm,
	).Scan(&reembolso.ID)

	if err != nil {
		return nil, err
	}
	return reembolso, nil
}

func (r *ReembolsoRepository) BuscarPorID(id int64) (*model.Reembolso, error) {
	query := `
		SELECT id, usuario_id, descricao, valor, categoria, status, comprovante, criado_em, atualizado_em
		FROM reembolsos WHERE id = $1
	`
	reembolso := &model.Reembolso{}
	err := r.db.QueryRow(query, id).Scan(
		&reembolso.ID,
		&reembolso.UsuarioID,
		&reembolso.Descricao,
		&reembolso.Valor,
		&reembolso.Categoria,
		&reembolso.Status,
		&reembolso.Comprovante,
		&reembolso.CriadoEm,
		&reembolso.AtualizadoEm,
	)
	if err != nil {
		return nil, err
	}
	return reembolso, nil
}

func (r *ReembolsoRepository) ListarPorUsuario(usuarioID int64) ([]model.Reembolso, error) {
	query := `
		SELECT id, usuario_id, descricao, valor, categoria, status, comprovante, criado_em, atualizado_em
		FROM reembolsos WHERE usuario_id = $1 ORDER BY criado_em DESC
	`
	rows, err := r.db.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reembolsos []model.Reembolso
	for rows.Next() {
		var reembolso model.Reembolso
		if err := rows.Scan(
			&reembolso.ID,
			&reembolso.UsuarioID,
			&reembolso.Descricao,
			&reembolso.Valor,
			&reembolso.Categoria,
			&reembolso.Status,
			&reembolso.Comprovante,
			&reembolso.CriadoEm,
			&reembolso.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		reembolsos = append(reembolsos, reembolso)
	}
	return reembolsos, nil
}

func (r *ReembolsoRepository) AtualizarStatus(id int64, status model.StatusReembolso) error {
	query := `UPDATE reembolsos SET status = $1, atualizado_em = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

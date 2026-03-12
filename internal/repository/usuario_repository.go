package repository

import (
	"database/sql"
	"time"

	"reembolso/internal/model"
)

type UsuarioRepository struct {
	db *sql.DB
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

type UsuarioComSenha struct {
	model.Usuario
	SenhaHash string
}

func (r *UsuarioRepository) Criar(nome, email, senhaHash string) (*model.Usuario, error) {
	query := `
		INSERT INTO usuarios (nome, email, senha_hash, criado_em)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	usuario := &model.Usuario{
		Nome:      nome,
		Email:     email,
		CriadoEm: time.Now(),
	}
	err := r.db.QueryRow(query, nome, email, senhaHash, usuario.CriadoEm).Scan(&usuario.ID)
	if err != nil {
		return nil, err
	}
	return usuario, nil
}

func (r *UsuarioRepository) BuscarPorEmail(email string) (*UsuarioComSenha, error) {
	query := `SELECT id, nome, email, senha_hash, criado_em FROM usuarios WHERE email = $1`
	u := &UsuarioComSenha{}
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Nome, &u.Email, &u.SenhaHash, &u.CriadoEm,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UsuarioRepository) BuscarPorID(id int64) (*model.Usuario, error) {
	query := `SELECT id, nome, email, criado_em FROM usuarios WHERE id = $1`
	u := &model.Usuario{}
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Nome, &u.Email, &u.CriadoEm)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UsuarioRepository) EmailExiste(email string) (bool, error) {
	var existe bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = $1)`, email).Scan(&existe)
	return existe, err
}

package seeds

import (
	"database/sql"
	"fmt"

	"reembolso/internal/auth"
	"reembolso/internal/logger"
)

// Seeder — equivalente à interface Seeder do Laravel
type Seeder interface {
	Run(db *sql.DB) error
	Name() string
}

// SeederRunner — equivalente ao DatabaseSeeder do Laravel
type SeederRunner struct {
	db      *sql.DB
	seeders []Seeder
}

func NewSeederRunner(db *sql.DB) *SeederRunner {
	return &SeederRunner{
		db: db,
		seeders: []Seeder{
			&UsuarioSeeder{},
			&ReembolsoSeeder{},
		},
	}
}

// Run — equivalente ao php artisan db:seed
func (r *SeederRunner) Run() error {
	for _, s := range r.seeders {
		logger.Info("rodando seeder", logger.With("seeder", s.Name()))
		if err := s.Run(r.db); err != nil {
			return fmt.Errorf("seeder %s falhou: %w", s.Name(), err)
		}
		logger.Info("seeder concluído", logger.With("seeder", s.Name()))
	}
	return nil
}

// ──────────────────────────────────────────────
// UsuarioSeeder
// ──────────────────────────────────────────────

type UsuarioSeeder struct{}

func (s *UsuarioSeeder) Name() string { return "UsuarioSeeder" }

func (s *UsuarioSeeder) Run(db *sql.DB) error {
	usuarios := []struct {
		nome  string
		email string
		senha string
	}{
		{"Admin", "admin@exemplo.com", "password123"},
		{"João Silva", "joao@exemplo.com", "password123"},
		{"Maria Souza", "maria@exemplo.com", "password123"},
	}

	for _, u := range usuarios {
		hash, err := auth.HashSenha(u.senha)
		if err != nil {
			return err
		}
		_, err = db.Exec(`
			INSERT INTO usuarios (nome, email, senha_hash)
			VALUES ($1, $2, $3)
			ON CONFLICT (email) DO NOTHING
		`, u.nome, u.email, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

// ──────────────────────────────────────────────
// ReembolsoSeeder
// ──────────────────────────────────────────────

type ReembolsoSeeder struct{}

func (s *ReembolsoSeeder) Name() string { return "ReembolsoSeeder" }

func (s *ReembolsoSeeder) Run(db *sql.DB) error {
	// Busca o primeiro usuário para associar os reembolsos
	var usuarioID int64
	if err := db.QueryRow(`SELECT id FROM usuarios LIMIT 1`).Scan(&usuarioID); err != nil {
		return fmt.Errorf("nenhum usuário encontrado para seeder de reembolsos: %w", err)
	}

	reembolsos := []struct {
		descricao string
		valor     float64
		categoria string
		status    string
	}{
		{"Passagem aérea para SP", 450.00, "TRANSPORTE", "PENDENTE"},
		{"Hotel - Conferência Tech", 320.00, "HOSPEDAGEM", "APROVADO"},
		{"Almoço com cliente", 89.90, "ALIMENTACAO", "PENDENTE"},
		{"Uber para reunião", 45.50, "TRANSPORTE", "REJEITADO"},
		{"Material de escritório", 120.00, "MATERIAL", "APROVADO"},
	}

	for _, r := range reembolsos {
		_, err := db.Exec(`
			INSERT INTO reembolsos (usuario_id, descricao, valor, categoria, status)
			VALUES ($1, $2, $3, $4, $5)
		`, usuarioID, r.descricao, r.valor, r.categoria, r.status)
		if err != nil {
			return err
		}
	}
	return nil
}

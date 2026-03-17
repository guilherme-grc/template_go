package seeds

import (
	"database/sql"
	"fmt"

	"template.go/internal/auth"
	"template.go/internal/logger"
)

// Seeder - Equivalent to Laravel's Seeder interface
type Seeder interface {
	Run(db *sql.DB) error
	Name() string
}

// SeederRunner - Equivalent to Laravel's DatabaseSeeder
type SeederRunner struct {
	db      *sql.DB
	seeders []Seeder
}

func NewSeederRunner(db *sql.DB) *SeederRunner {
	return &SeederRunner{
		db: db,
		seeders: []Seeder{
			&UserSeeder{},
		},
	}
}

// Run - Equivalent to 'php artisan db:seed'
func (r *SeederRunner) Run() error {
	for _, s := range r.seeders {
		logger.Info("running seeder", logger.With("seeder", s.Name()))

		if err := s.Run(r.db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", s.Name(), err)
		}

		logger.Info("seeder completed", logger.With("seeder", s.Name()))
	}
	return nil
}

// ---------------------------------------------------------
// UserSeeder
// ---------------------------------------------------------

type UserSeeder struct{}

func (s *UserSeeder) Name() string { return "UserSeeder" }

func (s *UserSeeder) Run(db *sql.DB) error {
	users := []struct {
		name     string
		email    string
		password string
	}{
		{"Admin", "admin@example.com", "password123"},
		{"John Doe", "john@example.com", "password123"},
		{"Jane Smith", "jane@example.com", "password123"},
	}

	for _, u := range users {
		// Note: Changed from HashSenha to HashPassword to follow English naming
		hash, err := auth.HashPassword(u.password)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			INSERT INTO users (name, email, password)
			VALUES ($1, $2, $3)
			ON CONFLICT (email) DO NOTHING
		`, u.name, u.email, hash)

		if err != nil {
			return err
		}
	}
	return nil
}

package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"reembolso/internal/logger"
)

// Migrator — equivalente ao sistema de migrations do Laravel
type Migrator struct {
	db         *sql.DB
	migrationsPath string
}

func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{db: db, migrationsPath: migrationsPath}
}

// Run — equivalente ao php artisan migrate
func (m *Migrator) Run() error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("erro ao criar tabela de migrations: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(m.migrationsPath, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files) // garante ordem (001, 002, 003...)

	for _, file := range files {
		name := filepath.Base(file)
		if m.jaRodou(name) {
			logger.Info("migration já aplicada, pulando", logger.With("migration", name))
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("erro ao ler migration %s: %w", name, err)
		}

		// Executa cada statement separadamente
		statements := splitStatements(string(content))
		for _, stmt := range statements {
			if strings.TrimSpace(stmt) == "" {
				continue
			}
			if _, err := m.db.Exec(stmt); err != nil {
				return fmt.Errorf("erro ao executar migration %s: %w", name, err)
			}
		}

		if err := m.registrar(name); err != nil {
			return err
		}
		logger.Info("migration aplicada", logger.With("migration", name))
	}

	logger.Info("migrations concluídas")
	return nil
}

// Fresh — equivalente ao php artisan migrate:fresh (drop + migrate)
func (m *Migrator) Fresh() error {
	logger.Warning("executando migrate:fresh — todos os dados serão apagados!")
	if _, err := m.db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`); err != nil {
		return err
	}
	return m.Run()
}

func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL UNIQUE,
			executado_em TIMESTAMP DEFAULT NOW()
		)
	`)
	return err
}

func (m *Migrator) jaRodou(name string) bool {
	var count int
	m.db.QueryRow(`SELECT COUNT(*) FROM migrations WHERE migration = $1`, name).Scan(&count)
	return count > 0
}

func (m *Migrator) registrar(name string) error {
	_, err := m.db.Exec(`INSERT INTO migrations (migration, executado_em) VALUES ($1, $2)`, name, time.Now())
	return err
}

func splitStatements(sql string) []string {
	return strings.Split(sql, ";")
}

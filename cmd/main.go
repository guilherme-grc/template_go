package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"reembolso/config"
	"reembolso/db"
	"reembolso/db/seeds"
	"reembolso/internal/auth"
	"reembolso/internal/handler"
	"reembolso/internal/logger"
	"reembolso/internal/repository"
	"reembolso/internal/service"
)

func main() {
	// Flags — equivalente aos comandos artisan do Laravel
	// go run cmd/main.go --migrate
	// go run cmd/main.go --seed
	// go run cmd/main.go --migrate --seed
	migrateFlag := flag.Bool("migrate", false, "Rodar migrations (php artisan migrate)")
	seedFlag := flag.Bool("seed", false, "Rodar seeders (php artisan db:seed)")
	freshFlag := flag.Bool("fresh", false, "Drop + migrate (php artisan migrate:fresh)")
	flag.Parse()

	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Banco inacessível: %v", err)
	}
	logger.Info("conectado ao banco de dados")

	// ── Migrations ──────────────────────────────────────────────────────────
	migrator := db.NewMigrator(sqlDB, "db/migrations")

	if *freshFlag {
		if err := migrator.Fresh(); err != nil {
			log.Fatalf("migrate:fresh falhou: %v", err)
		}
	} else if *migrateFlag {
		if err := migrator.Run(); err != nil {
			log.Fatalf("migrate falhou: %v", err)
		}
	}

	// ── Seeders ─────────────────────────────────────────────────────────────
	if *seedFlag {
		runner := seeds.NewSeederRunner(sqlDB)
		if err := runner.Run(); err != nil {
			log.Fatalf("seeder falhou: %v", err)
		}
		if !*migrateFlag && !*freshFlag {
			os.Exit(0) // só seed, sem servidor
		}
	}

	// ── Injeção de dependências ──────────────────────────────────────────────
	jwtSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTAccessExpiryMin, cfg.JWTRefreshExpiryDays)

	usuarioRepo := repository.NewUsuarioRepository(sqlDB)
	reembolsoRepo := repository.NewReembolsoRepository(sqlDB)

	authSvc := service.NewAuthService(usuarioRepo, jwtSvc)
	reembolsoSvc := service.NewReembolsoService(reembolsoRepo)

	// ── Servidor ─────────────────────────────────────────────────────────────
	mux := http.NewServeMux()
	appHandler := handler.RegisterRoutes(mux, authSvc, reembolsoSvc, jwtSvc)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      appHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful Shutdown ────────────────────────────────────────────────────
	// Equivalente ao SIGTERM handling do Laravel Octane / supervisord
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("servidor iniciado", logger.With("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	<-done // aguarda sinal de encerramento
	logger.Warning("encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro no shutdown: %v", err)
	}
	logger.Info("servidor encerrado com sucesso")
}

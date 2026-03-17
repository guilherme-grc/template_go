#Template API
A robust REST API built in Go for managing corporate reimbursements, featuring JWT authentication, structured event logging, and a modular architecture inspired by Laravel’s clear separation of concerns.

## Estrutura

```
template-api/
├── cmd/
│               # Application entry point (main.go)
├── internal/
│   ├── auth/            # JWT (Access/Refresh), Bcrypt, and Context logic
│   ├── config/          # Environment variable loading (godotenv)
│   ├── handler/         # HTTP Handlers (Controllers)
│   ├── logger/          # Structured JSON logging (slog)
│   ├── middleware/      # Auth, CORS, Rate Limiting, and Recovery
│   ├── model/           # Entities (User, Reimbursement) and Request DTOs
│   ├── repository/      # Database Access Layer (PostgreSQL)
│   ├── service/         # Business Logic Layer
│   └── validation/      # Custom validation logic
├── seeds/               # Database seeders
└── .env                 # Environment variables (JWT_SECRET, DB_URL, etc.)
```

## Como rodar

```bash
# Install dependencies
go mod tidy

# Run the application (Server only)
go run ./cmd

# Database Management (CLI Commands)
go run ./cmd --migrate           # Run migrations
go run ./cmd --migrate --seed    # Run migrations and seeders
go run ./cmd --fresh --seed      # Drop all tables, migrate, and seed
```

## Endpoints

| Método | Rota                        | Descrição              |
|--------|-----------------------------|------------------------|
| GET    | /health                     | Health check           |



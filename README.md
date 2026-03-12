# Sistema de Reembolso

API REST em Go para gerenciamento de reembolsos.

## Estrutura

```
reembolso/
├── cmd/           # Ponto de entrada da aplicação
├── config/        # Configurações e variáveis de ambiente
├── internal/
│   ├── handler/   # HTTP handlers
│   ├── service/   # Regras de negócio
│   ├── repository/# Acesso ao banco de dados
│   └── model/     # Entidades e structs
└── db/migrations/ # Scripts SQL
```

## Como rodar

```bash
# Instalar dependências
go mod tidy

# Rodar a aplicação
go run cmd/main.go
```

## Endpoints

| Método | Rota                        | Descrição              |
|--------|-----------------------------|------------------------|
| GET    | /health                     | Health check           |
| POST   | /reembolsos                 | Criar reembolso        |
| GET    | /reembolsos/{id}            | Buscar por ID          |
| GET    | /reembolsos?usuario_id=1    | Listar por usuário     |
| PATCH  | /reembolsos/{id}/aprovar    | Aprovar reembolso      |
| PATCH  | /reembolsos/{id}/rejeitar   | Rejeitar reembolso     |

## Status possíveis

- `PENDENTE` → estado inicial
- `APROVADO` → reembolso aprovado
- `REJEITADO` → reembolso rejeitado

```bash

go mod tidy        # baixa jwt, bcrypt, pq
go run cmd/main.go

```
```

**Fluxo de uso:**

POST /auth/register  → retorna access_token + refresh_token
POST /auth/login     → retorna access_token + refresh_token
POST /auth/refresh   → renova com o refresh_token
GET  /auth/me        → Authorization: Bearer <access_token>
POST /reembolsos     → Authorization: Bearer <access_token>

go mod tidy                              # instalar dependências
go run cmd/main.go --migrate             # rodar migrations
go run cmd/main.go --migrate --seed      # migrations + seeders
go run cmd/main.go --fresh --seed        # drop + migrate + seed
go run cmd/main.go                       # só sobe o servidor
```
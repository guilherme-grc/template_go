package model

import "time"

type StatusReembolso string

const (
	StatusPendente  StatusReembolso = "PENDENTE"
	StatusAprovado  StatusReembolso = "APROVADO"
	StatusRejeitado StatusReembolso = "REJEITADO"
)

type Reembolso struct {
	ID           int64           `json:"id"`
	UsuarioID    int64           `json:"usuario_id"`
	Descricao    string          `json:"descricao"`
	Valor        float64         `json:"valor"`
	Categoria    string          `json:"categoria"`
	Status       StatusReembolso `json:"status"`
	Comprovante  string          `json:"comprovante,omitempty"`
	CriadoEm    time.Time       `json:"criado_em"`
	AtualizadoEm time.Time      `json:"atualizado_em"`
}

type CriarReembolsoRequest struct {
	UsuarioID   int64   `json:"usuario_id"`
	Descricao   string  `json:"descricao"`
	Valor       float64 `json:"valor"`
	Categoria   string  `json:"categoria"`
	Comprovante string  `json:"comprovante,omitempty"`
}

type AtualizarStatusRequest struct {
	Status StatusReembolso `json:"status"`
}

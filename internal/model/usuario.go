package model

import "time"

type Usuario struct {
	ID        int64     `json:"id"`
	Nome      string    `json:"nome"`
	Email     string    `json:"email"`
	CriadoEm time.Time `json:"criado_em"`
}

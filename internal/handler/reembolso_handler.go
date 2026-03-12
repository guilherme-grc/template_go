package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"reembolso/internal/auth"
	"reembolso/internal/model"
	"reembolso/internal/service"
)

type ReembolsoHandler struct {
	service *service.ReembolsoService
}

func NewReembolsoHandler(svc *service.ReembolsoService) *ReembolsoHandler {
	return &ReembolsoHandler{service: svc}
}

// Roteador distribui /reembolsos/{id} e /reembolsos/{id}/acao
func (h *ReembolsoHandler) Roteador(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "rota inválida", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}
	if len(parts) == 3 {
		switch parts[2] {
		case "aprovar":
			h.aprovar(w, r, id)
		case "rejeitar":
			h.rejeitar(w, r, id)
		default:
			http.Error(w, "ação desconhecida", http.StatusNotFound)
		}
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.buscar(w, r, id)
	default:
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
	}
}

// POST /reembolsos
func (h *ReembolsoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := auth.UsuarioDoContexto(r.Context())
	if !ok {
		http.Error(w, "não autenticado", http.StatusUnauthorized)
		return
	}
	var req model.CriarReembolsoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "body inválido", http.StatusBadRequest)
		return
	}
	req.UsuarioID = claims.UsuarioID
	reembolso, err := h.service.CriarReembolso(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reembolso)
}

func (h *ReembolsoHandler) buscar(w http.ResponseWriter, r *http.Request, id int64) {
	reembolso, err := h.service.BuscarReembolso(id)
	if err != nil {
		http.Error(w, "reembolso não encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reembolso)
}

func (h *ReembolsoHandler) aprovar(w http.ResponseWriter, r *http.Request, id int64) {
	if r.Method != http.MethodPatch {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := h.service.AprovarReembolso(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ReembolsoHandler) rejeitar(w http.ResponseWriter, r *http.Request, id int64) {
	if r.Method != http.MethodPatch {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := h.service.RejeitarReembolso(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

package http

import (
	"encoding/json"
	"net/http"
	"github.com/akozie/babe-25th-backend/internal/domain"
)

type GuestbookHandler struct {
	Usecase domain.GuestbookUsecase
}

func (h *GuestbookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var entry domain.GuestbookEntry
	json.NewDecoder(r.Body).Decode(&entry)
	h.Usecase.CreateEntry(r.Context(), &entry)
	w.WriteHeader(http.StatusCreated)
}

func (h *GuestbookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	entries, _ := h.Usecase.GetAllEntries(r.Context())
	json.NewEncoder(w).Encode(entries)
}
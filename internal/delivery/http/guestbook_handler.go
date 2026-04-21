package http

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

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
	log.Info("GET /api/v1/guestbook: handler entered")
	entries, err := h.Usecase.GetAllEntries(r.Context())
	if err != nil {
		log.WithError(err).Warn("GET /api/v1/guestbook failed")
		http.Error(w, "Could not fetch guestbook entries", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(entries)
}

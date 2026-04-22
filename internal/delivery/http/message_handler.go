package http

import (
	"encoding/json"
	"net/http"

	"github.com/akozie/babe-25th-backend/internal/domain"
)

type MessageHandler struct {
	Usecase domain.MessageUsecase // Assuming you have a usecase interface
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var msg domain.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Save to database
	err := h.Usecase.CreateMessage(r.Context(), &msg)
	if err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (h *MessageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	messages, err := h.Usecase.GetAllMessages(r.Context())
	if err != nil {
		writeServiceError(w, err, "GET /api/v1/messages failed")
		return
	}
	json.NewEncoder(w).Encode(messages)
}

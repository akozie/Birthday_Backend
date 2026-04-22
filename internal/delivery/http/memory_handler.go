package http

import (
	"encoding/json"
	"net/http"

	"github.com/akozie/babe-25th-backend/internal/domain"
)

type MemoryHandler struct {
	Usecase domain.MemoryUsecase
}

func (h *MemoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the multipart form (for the file)
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Create the memory object from form values
	memory := &domain.Memory{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	// 3. Call the Usecase
	err = h.Usecase.CreateMemory(r.Context(), memory, file)
	if err != nil {
		http.Error(w, "Failed to save memory", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(memory)
}

func (h *MemoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	memories, err := h.Usecase.GetAllMemories(r.Context())
	if err != nil {
		writeServiceError(w, err, "GET /api/v1/memories failed")
		return
	}
	json.NewEncoder(w).Encode(memories)
}

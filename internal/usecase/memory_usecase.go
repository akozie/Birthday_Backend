package usecase

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/akozie/babe-25th-backend/internal/domain"
	"github.com/akozie/babe-25th-backend/pkg/media"
)

type memoryUsecase struct {
	repo           domain.MemoryRepository
	cloudinary     *media.CloudinaryService
	contextTimeout time.Duration
}

func NewMemoryUsecase(r domain.MemoryRepository, c *media.CloudinaryService, timeout time.Duration) domain.MemoryUsecase {
	return &memoryUsecase{repo: r, cloudinary: c, contextTimeout: timeout}
}

func (u *memoryUsecase) CreateMemory(ctx context.Context, memory *domain.Memory, file multipart.File) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1. Logic: Upload to Cloudinary first
	url, err := u.cloudinary.UploadFile(ctx, file)
	if err != nil {
		return err
	}

	// 2. Logic: Update the memory struct with the new URL
	memory.MediaURL = url

	// 3. Logic: Save to MongoDB
	return u.repo.Create(ctx, memory)
}

func (u *memoryUsecase) GetAllMemories(ctx context.Context) ([]domain.Memory, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.FetchAll(ctx)
}

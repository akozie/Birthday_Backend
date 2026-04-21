package usecase

import (
	"context"
	"mime/multipart"
	
	"github.com/akozie/babe-25th-backend/internal/domain"
	"github.com/akozie/babe-25th-backend/pkg/media"
)

type memoryUsecase struct {
	repo         domain.MemoryRepository
	cloudinary   *media.CloudinaryService
}

func NewMemoryUsecase(r domain.MemoryRepository, c *media.CloudinaryService) domain.MemoryUsecase {
	return &memoryUsecase{repo: r, cloudinary: c}
}

func (u *memoryUsecase) CreateMemory(ctx context.Context, memory *domain.Memory, file multipart.File) error {
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
	return u.repo.FetchAll(ctx)
}
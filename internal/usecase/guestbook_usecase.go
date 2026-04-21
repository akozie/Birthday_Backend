package usecase

import (
	"context"
	"github.com/akozie/babe-25th-backend/internal/domain"
)

type guestbookUsecase struct {
	repo domain.GuestbookRepository
}

func NewGuestbookUsecase(r domain.GuestbookRepository) domain.GuestbookUsecase {
	return &guestbookUsecase{repo: r}
}

func (u *guestbookUsecase) CreateEntry(ctx context.Context, entry *domain.GuestbookEntry) error {
	return u.repo.Create(ctx, entry)
}

func (u *guestbookUsecase) GetAllEntries(ctx context.Context) ([]domain.GuestbookEntry, error) {
	return u.repo.FetchAll(ctx)
}
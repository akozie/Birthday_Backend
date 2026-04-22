package usecase

import (
	"context"
	"time"

	"github.com/akozie/babe-25th-backend/internal/domain"
)

type guestbookUsecase struct {
	repo           domain.GuestbookRepository
	contextTimeout time.Duration
}

func NewGuestbookUsecase(r domain.GuestbookRepository, timeout time.Duration) domain.GuestbookUsecase {
	return &guestbookUsecase{repo: r, contextTimeout: timeout}
}

func (u *guestbookUsecase) CreateEntry(ctx context.Context, entry *domain.GuestbookEntry) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.Create(ctx, entry)
}

func (u *guestbookUsecase) GetAllEntries(ctx context.Context) ([]domain.GuestbookEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.FetchAll(ctx)
}

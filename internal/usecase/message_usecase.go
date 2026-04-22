package usecase

import (
	"context"
	"time"

	"github.com/akozie/babe-25th-backend/internal/domain"
)

type messageUsecase struct {
	repo           domain.MessageRepository
	contextTimeout time.Duration
}

// NewMessageUsecase is the constructor
func NewMessageUsecase(repo domain.MessageRepository, timeout time.Duration) domain.MessageUsecase {
	return &messageUsecase{repo: repo, contextTimeout: timeout}
}

func (u *messageUsecase) CreateMessage(ctx context.Context, m *domain.Message) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	m.CreatedAt = time.Now() 
	return u.repo.Create(ctx, m)
}

// Added the missing method required by your interface
func (u *messageUsecase) GetAllMessages(ctx context.Context) ([]domain.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.FetchAll(ctx)
}

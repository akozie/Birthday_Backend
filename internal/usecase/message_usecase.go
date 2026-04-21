package usecase

import (
	"context"
	"time"
	"github.com/akozie/babe-25th-backend/internal/domain"
)

type messageUsecase struct {
	repo domain.MessageRepository
}

// NewMessageUsecase is the constructor
func NewMessageUsecase(repo domain.MessageRepository) domain.MessageUsecase {
	return &messageUsecase{repo: repo}
}

func (u *messageUsecase) CreateMessage(ctx context.Context, m *domain.Message) error {
	m.CreatedAt = time.Now() 
	return u.repo.Create(ctx, m)
}

// Added the missing method required by your interface
func (u *messageUsecase) GetAllMessages(ctx context.Context) ([]domain.Message, error) {
	return u.repo.FetchAll(ctx)
}
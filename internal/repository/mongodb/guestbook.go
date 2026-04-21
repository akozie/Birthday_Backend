package mongodb

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/akozie/babe-25th-backend/internal/domain"
)

type mongoGuestbookRepo struct {
	db *mongo.Collection
}

func NewMongoGuestbookRepository(db *mongo.Database) domain.GuestbookRepository {
	return &mongoGuestbookRepo{db: db.Collection("guestbook")}
}

func (r *mongoGuestbookRepo) Create(ctx context.Context, entry *domain.GuestbookEntry) error {
	entry.CreatedAt = time.Now()
	_, err := r.db.InsertOne(ctx, entry)
	return err
}

func (r *mongoGuestbookRepo) FetchAll(ctx context.Context) ([]domain.GuestbookEntry, error) {
	// Initialize as empty slice to avoid 'null' response
	entries := []domain.GuestbookEntry{}
	cursor, err := r.db.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
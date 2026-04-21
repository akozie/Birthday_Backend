package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/akozie/babe-25th-backend/internal/domain"
)

type messageRepo struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) domain.MessageRepository {
	return &messageRepo{collection: db.Collection("messages")}
}

func (r *messageRepo) Create(ctx context.Context, m *domain.Message) error {
	_, err := r.collection.InsertOne(ctx, m)
	return err
}

func (r *messageRepo) FetchAll(ctx context.Context) ([]domain.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil { return nil, err }
	var messages []domain.Message
	err = cursor.All(ctx, &messages)
	return messages, err
}
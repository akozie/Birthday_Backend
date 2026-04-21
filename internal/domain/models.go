package domain

import (
	"context"
	"time"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Memory represents a milestone or moment in your relationship/her life
type Memory struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	MediaURL     string             `json:"media_url" bson:"media_url"`     // From Cloudinary
	MediaType    string             `json:"media_type" bson:"media_type"`   // e.g., "image" or "video"
	DateOccurred time.Time          `json:"date_occurred" bson:"date_occurred"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

// GuestbookEntry represents a birthday wish from a friend or family member
type GuestbookEntry struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Author    string             `json:"author" bson:"author"`
	Message   string             `json:"message" bson:"message"`
	MediaURL  string             `json:"media_url" bson:"media_url"` // From Cloudinary (optional)
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Message struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}


// MemoryRepository defines how we interact with Memory data
type MemoryRepository interface {
	Create(ctx context.Context, memory *Memory) error
	FetchAll(ctx context.Context) ([]Memory, error)
	GetByID(ctx context.Context, id string) (*Memory, error)
}

// GuestbookRepository defines how we interact with Guestbook data
type GuestbookRepository interface {
	Create(ctx context.Context, entry *GuestbookEntry) error
	FetchAll(ctx context.Context) ([]GuestbookEntry, error)
}

// MessageRepository defines the methods that the database layer must implement
type MessageRepository interface {
	Create(ctx context.Context, m *Message) error
	FetchAll(ctx context.Context) ([]Message, error)
}

type MemoryUsecase interface {
	CreateMemory(ctx context.Context, memory *Memory, file multipart.File) error
	GetAllMemories(ctx context.Context) ([]Memory, error)
}

// GuestbookUsecase defines the business logic rules for the guestbook
type GuestbookUsecase interface {
	CreateEntry(ctx context.Context, entry *GuestbookEntry) error
	GetAllEntries(ctx context.Context) ([]GuestbookEntry, error)
}

// MessageUsecase defines the business logic layer
type MessageUsecase interface {
	CreateMessage(ctx context.Context, m *Message) error
	GetAllMessages(ctx context.Context) ([]Message, error)
}
package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// Replace this with your actual module name from go.mod
	"github.com/akozie/babe-25th-backend/internal/domain" 
)

// mongoMemoryRepo is the concrete implementation of domain.MemoryRepository
type mongoMemoryRepo struct {
	db *mongo.Collection
}

// NewMongoMemoryRepository creates a new MongoDB repository for memories
func NewMongoMemoryRepository(db *mongo.Database) domain.MemoryRepository {
	return &mongoMemoryRepo{
		db: db.Collection("memories"),
	}
}

// Create inserts a new memory into the MongoDB collection
func (m *mongoMemoryRepo) Create(ctx context.Context, memory *domain.Memory) error {
	// Set the creation timestamp if it's empty
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}

	result, err := m.db.InsertOne(ctx, memory)
	if err != nil {
		return err
	}

	// Attach the newly generated MongoDB ObjectID back to our struct
	memory.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FetchAll retrieves all memories, sorted by the date they occurred (oldest to newest for a timeline)
func (m *mongoMemoryRepo) FetchAll(ctx context.Context) ([]domain.Memory, error) {
	//var memories []domain.Memory
	memories := []domain.Memory{}

	// Sort by DateOccurred ascending (1)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "date_occurred", Value: 1}})

	cursor, err := m.db.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the MongoDB BSON documents into our clean Go structs
	for cursor.Next(ctx) {
		var memory domain.Memory
		if err := cursor.Decode(&memory); err != nil {
			return nil, err
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// GetByID fetches a single memory using its unique ID
func (m *mongoMemoryRepo) GetByID(ctx context.Context, id string) (*domain.Memory, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var memory domain.Memory
	err = m.db.FindOne(ctx, bson.M{"_id": objectID}).Decode(&memory)
	if err != nil {
		return nil, err
	}

	return &memory, nil
}
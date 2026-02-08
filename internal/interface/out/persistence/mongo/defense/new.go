package defense

import (
	domain "counter-scam-agency/internal/domain/defense"

	"go.mongodb.org/mongo-driver/mongo"
)

// Compile-time interface guard.
var _ domain.BaseRepository = (*MongoRepository)(nil)

// MongoRepository stores defense base data in MongoDB.
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a MongoDB-backed defense base repository.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("bases"),
	}
}

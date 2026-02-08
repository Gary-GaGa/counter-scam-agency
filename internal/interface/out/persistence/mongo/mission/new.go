package mission

import (
	"counter-scam-agency/internal/domain/intelligence"

	"go.mongodb.org/mongo-driver/mongo"
)

// Compile-time interface guard.
var _ intelligence.MissionRepository = (*MongoRepository)(nil)

// MongoRepository stores mission data in MongoDB.
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a MongoDB-backed mission repository.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("missions"),
	}
}

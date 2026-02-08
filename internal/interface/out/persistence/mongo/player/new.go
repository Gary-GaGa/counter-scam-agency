package player

import (
	"counter-scam-agency/internal/domain/personnel"

	"go.mongodb.org/mongo-driver/mongo"
)

// Compile-time interface guard.
var _ personnel.PlayerRepository = (*MongoRepository)(nil)

// MongoRepository stores player data in MongoDB.
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a MongoDB-backed player repository.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("players"),
	}
}

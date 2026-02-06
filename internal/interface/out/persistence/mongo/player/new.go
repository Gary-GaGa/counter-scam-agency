package player

import "go.mongodb.org/mongo-driver/mongo"

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

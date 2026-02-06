package mission

import "go.mongodb.org/mongo-driver/mongo"

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

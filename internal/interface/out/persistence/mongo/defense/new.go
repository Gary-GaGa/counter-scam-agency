package defense

import "go.mongodb.org/mongo-driver/mongo"

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

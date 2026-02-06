package investigation

import "go.mongodb.org/mongo-driver/mongo"

// MongoRepository stores investigation data in MongoDB.
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a MongoDB-backed investigation repository.
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("investigations"),
	}
}

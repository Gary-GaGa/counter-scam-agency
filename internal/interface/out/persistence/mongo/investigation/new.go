package investigation

import (
	"counter-scam-agency/internal/domain/operation"

	"go.mongodb.org/mongo-driver/mongo"
)

// Compile-time interface guard.
var _ operation.InvestigationRepository = (*MongoRepository)(nil)

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

package defense

import "go.mongodb.org/mongo-driver/mongo"

// defenseRepo -
type defenseRepo struct {
	collection *mongo.Collection
}

// NewDefenseRepo -
func NewDefenseRepo(db *mongo.Database) *defenseRepo {
	return &defenseRepo{
		collection: db.Collection("bases"),
	}
}

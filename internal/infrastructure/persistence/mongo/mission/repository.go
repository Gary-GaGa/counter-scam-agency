package mission

import (
	"context"

	"counter-scam-agency/internal/domain/intelligence"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("missions"),
	}
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*intelligence.Mission, error) {
	var mission intelligence.Mission
	filter := bson.M{"id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&mission)
	if err != nil {
		return nil, err
	}
	return &mission, nil
}

func (r *MongoRepository) FindAll(ctx context.Context) ([]*intelligence.Mission, error) {
	var missions []*intelligence.Mission
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &missions); err != nil {
		return nil, err
	}
	return missions, nil
}

func (r *MongoRepository) Save(ctx context.Context, mission *intelligence.Mission) error {
	filter := bson.M{"id": mission.ID}
	update := bson.M{"$set": mission}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

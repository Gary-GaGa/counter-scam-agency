package mission

import (
	"context"
	"fmt"

	domain "counter-scam-agency/internal/domain/intelligence"
	modelconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/model"
	po "counter-scam-agency/internal/interface/out/persistence/mongo/po/intelligence"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindByID retrieves a mission by ID.
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*domain.Mission, error) {
	if id == "" {
		return nil, fmt.Errorf("mission id is required")
	}

	var doc po.MissionDocPo
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find mission: %w", err)
	}

	return modelconv.MissionDocToModel(&doc), nil
}

// FindAll retrieves all missions.
func (r *MongoRepository) FindAll(ctx context.Context) ([]*domain.Mission, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("find missions: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	results := make([]*domain.Mission, 0)
	for cursor.Next(ctx) {
		var doc po.MissionDocPo
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode mission: %w", err)
		}
		results = append(results, modelconv.MissionDocToModel(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate missions: %w", err)
	}

	return results, nil
}

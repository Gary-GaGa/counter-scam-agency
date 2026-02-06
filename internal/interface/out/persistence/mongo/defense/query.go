package defense

import (
	"context"
	"fmt"

	domain "counter-scam-agency/internal/domain/defense"
	modelconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/model"
	po "counter-scam-agency/internal/interface/out/persistence/mongo/po/defense"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindByID retrieves a base by ID.
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*domain.Base, error) {
	if id == "" {
		return nil, fmt.Errorf("base id is required")
	}

	var doc po.BaseDocPo
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find base: %w", err)
	}

	return modelconv.BaseDocToModel(&doc), nil
}

// FindByOwnerID retrieves bases by owner ID.
func (r *MongoRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*domain.Base, error) {
	if ownerID == "" {
		return nil, fmt.Errorf("owner id is required")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, fmt.Errorf("find bases: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	results := make([]*domain.Base, 0)
	for cursor.Next(ctx) {
		var doc po.BaseDocPo
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode base: %w", err)
		}
		results = append(results, modelconv.BaseDocToModel(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate bases: %w", err)
	}

	return results, nil
}

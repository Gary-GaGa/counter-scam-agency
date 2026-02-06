package investigation

import (
	"context"
	"fmt"

	domain "counter-scam-agency/internal/domain/operation"
	modelconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/model"
	po "counter-scam-agency/internal/interface/out/persistence/mongo/po/operation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindByID retrieves an investigation by ID.
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*domain.Investigation, error) {
	if id == "" {
		return nil, fmt.Errorf("investigation id is required")
	}

	var doc po.InvestigationDocPo
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find investigation: %w", err)
	}

	return modelconv.InvestigationDocToModel(&doc), nil
}

// FindByPlayerID retrieves investigations by player ID.
func (r *MongoRepository) FindByPlayerID(ctx context.Context, playerID string) ([]*domain.Investigation, error) {
	if playerID == "" {
		return nil, fmt.Errorf("player id is required")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"player_id": playerID})
	if err != nil {
		return nil, fmt.Errorf("find investigations: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	results := make([]*domain.Investigation, 0)
	for cursor.Next(ctx) {
		var doc po.InvestigationDocPo
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode investigation: %w", err)
		}
		results = append(results, modelconv.InvestigationDocToModel(&doc))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate investigations: %w", err)
	}

	return results, nil
}

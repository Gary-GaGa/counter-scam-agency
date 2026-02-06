package player

import (
	"context"
	"fmt"

	domain "counter-scam-agency/internal/domain/personnel"
	modelconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/model"
	po "counter-scam-agency/internal/interface/out/persistence/mongo/po/personnel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindByID retrieves a player by ID.
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*domain.Player, error) {
	if id == "" {
		return nil, fmt.Errorf("player id is required")
	}

	var doc po.PlayerDocPo
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find player: %w", err)
	}

	return modelconv.PlayerDocToModel(&doc), nil
}

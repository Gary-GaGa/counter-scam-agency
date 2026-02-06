package mission

import (
	"context"
	"fmt"

	"counter-scam-agency/internal/domain/intelligence"
	poconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/po"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save upserts a mission.
func (r *MongoRepository) Save(ctx context.Context, mission *intelligence.Mission) error {
	if mission == nil {
		return fmt.Errorf("mission is nil")
	}
	if mission.ID == "" {
		return fmt.Errorf("mission id is required")
	}

	doc := poconv.MissionDocToPo(mission)
	filter := bson.M{"id": mission.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("save mission: %w", err)
	}

	return nil
}

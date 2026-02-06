package player

import (
	"context"
	"fmt"

	"counter-scam-agency/internal/domain/personnel"
	poconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/po"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save upserts a player.
func (r *MongoRepository) Save(ctx context.Context, player *personnel.Player) error {
	if player == nil {
		return fmt.Errorf("player is nil")
	}
	if player.ID == "" {
		return fmt.Errorf("player id is required")
	}

	doc := poconv.PlayerDocToPo(player)
	filter := bson.M{"id": player.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("save player: %w", err)
	}

	return nil
}

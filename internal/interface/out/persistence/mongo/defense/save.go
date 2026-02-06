package defense

import (
	"context"
	"fmt"

	"counter-scam-agency/internal/domain/defense"
	poconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/po"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save upserts a defense base.
func (r *MongoRepository) Save(ctx context.Context, base *defense.Base) error {
	if base == nil {
		return fmt.Errorf("base is nil")
	}
	if base.ID == "" {
		return fmt.Errorf("base id is required")
	}

	doc := poconv.BaseDocToPo(base)
	filter := bson.M{"id": base.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("save base: %w", err)
	}

	return nil
}

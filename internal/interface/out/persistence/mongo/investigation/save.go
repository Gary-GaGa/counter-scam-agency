package investigation

import (
	"context"
	"fmt"

	"counter-scam-agency/internal/domain/operation"
	poconv "counter-scam-agency/internal/interface/out/persistence/mongo/po/convert/po"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save upserts an investigation.
func (r *MongoRepository) Save(ctx context.Context, inv *operation.Investigation) error {
	if inv == nil {
		return fmt.Errorf("investigation is nil")
	}
	if inv.ID == "" {
		return fmt.Errorf("investigation id is required")
	}

	doc := poconv.InvestigationDocToPo(inv)
	filter := bson.M{"id": inv.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("save investigation: %w", err)
	}

	return nil
}

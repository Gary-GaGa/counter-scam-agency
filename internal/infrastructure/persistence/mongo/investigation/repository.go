package investigation

import (
	"context"

	"counter-scam-agency/internal/domain/operation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("investigations"),
	}
}

type investigationDoc struct {
	ID                   string                        `bson:"id"`
	PlayerID             string                        `bson:"player_id"`
	MissionID            string                        `bson:"mission_id"`
	Status               operation.InvestigationStatus `bson:"status"`
	CurrentNodeID        string                        `bson:"current_node_id"`
	NodeHistory          []operation.NodeDecision      `bson:"node_history"`
	CollectedEvidenceIDs []string                      `bson:"collected_evidence_ids"`
	SuspicionLevel       int                           `bson:"suspicion_level"`
}

func toDoc(inv *operation.Investigation) *investigationDoc {
	return &investigationDoc{
		ID:                   inv.ID,
		PlayerID:             inv.PlayerID,
		MissionID:            inv.MissionID,
		Status:               inv.Status,
		CurrentNodeID:        inv.CurrentNodeID,
		NodeHistory:          inv.NodeHistory,
		CollectedEvidenceIDs: inv.CollectedEvidenceIDs,
		SuspicionLevel:       inv.SuspicionLevel,
	}
}

func fromDoc(doc *investigationDoc) *operation.Investigation {
	return &operation.Investigation{
		ID:                   doc.ID,
		PlayerID:             doc.PlayerID,
		MissionID:            doc.MissionID,
		Status:               doc.Status,
		CurrentNodeID:        doc.CurrentNodeID,
		NodeHistory:          doc.NodeHistory,
		CollectedEvidenceIDs: doc.CollectedEvidenceIDs,
		SuspicionLevel:       doc.SuspicionLevel,
	}
}

func (r *MongoRepository) Save(ctx context.Context, inv *operation.Investigation) error {
	filter := bson.M{"id": inv.ID}
	doc := toDoc(inv)
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*operation.Investigation, error) {
	var doc investigationDoc
	filter := bson.M{"id": id}
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}
	return fromDoc(&doc), nil
}

func (r *MongoRepository) FindByPlayerID(ctx context.Context, playerID string) ([]*operation.Investigation, error) {
	var docs []*investigationDoc
	filter := bson.M{"player_id": playerID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	invs := make([]*operation.Investigation, 0, len(docs))
	for _, doc := range docs {
		invs = append(invs, fromDoc(doc))
	}
	return invs, nil
}

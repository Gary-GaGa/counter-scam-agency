package player

import (
	"context"

	"counter-scam-agency/internal/domain/personnel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("players"),
	}
}

type playerDoc struct {
	ID                string               `bson:"id"`
	Stats             personnel.Stats      `bson:"stats"`
	Partner           *personnel.AIPartner `bson:"partner"`
	Reputation        int                  `bson:"reputation"`
	UnlockedModuleIDs []string             `bson:"unlocked_module_ids"`
}

func toDoc(player *personnel.Player) *playerDoc {
	unlocked := make([]string, 0, len(player.UnlockedModuleIDs))
	for id := range player.UnlockedModuleIDs {
		unlocked = append(unlocked, id)
	}
	return &playerDoc{
		ID:                player.ID,
		Stats:             player.Stats,
		Partner:           player.Partner,
		Reputation:        player.Reputation,
		UnlockedModuleIDs: unlocked,
	}
}

func fromDoc(doc *playerDoc) *personnel.Player {
	player := &personnel.Player{
		ID:                doc.ID,
		Stats:             doc.Stats,
		Partner:           doc.Partner,
		Reputation:        doc.Reputation,
		UnlockedModuleIDs: make(map[string]struct{}),
	}
	for _, id := range doc.UnlockedModuleIDs {
		player.UnlockedModuleIDs[id] = struct{}{}
	}
	return player
}

func (r *MongoRepository) Save(ctx context.Context, player *personnel.Player) error {
	filter := bson.M{"id": player.ID}
	doc := toDoc(player)
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*personnel.Player, error) {
	var doc playerDoc
	filter := bson.M{"id": id}
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}
	return fromDoc(&doc), nil
}

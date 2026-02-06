package operation

// InvestigationDocPo stores investigation data in MongoDB.
type InvestigationDocPo struct {
	ID                   string              `bson:"id"`
	PlayerID             string              `bson:"player_id"`
	MissionID            string              `bson:"mission_id"`
	Status               string              `bson:"status"`
	CurrentNodeID        string              `bson:"current_node_id"`
	NodeHistory          []NodeDecisionDocPo `bson:"node_history"`
	CollectedEvidenceIDs []string            `bson:"collected_evidence_ids"`
	SuspicionLevel       int                 `bson:"suspicion_level"`
}

// NodeDecisionDocPo stores node decisions.
type NodeDecisionDocPo struct {
	NodeID     string `bson:"node_id"`
	OptionID   string `bson:"option_id"`
	NextNodeID string `bson:"next_node_id"`
}

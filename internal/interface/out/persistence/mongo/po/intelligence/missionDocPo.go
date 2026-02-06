package intelligence

// MissionDocPo stores mission data in MongoDB.
type MissionDocPo struct {
	ID               string               `bson:"id"`
	Title            string               `bson:"title"`
	Description      string               `bson:"description"`
	Type             string               `bson:"type"`
	Difficulty       int                  `bson:"difficulty"`
	ReputationWeight int                  `bson:"reputation_weight"`
	VictimProfile    *VictimProfileDocPo  `bson:"victim_profile"`
	Nodes            []NarrativeNodeDocPo `bson:"nodes"`
	EvidenceList     []EvidenceDocPo      `bson:"evidence_list"`
}

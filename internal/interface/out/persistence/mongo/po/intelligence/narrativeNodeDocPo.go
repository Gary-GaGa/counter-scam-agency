package intelligence

// NarrativeNodeDocPo stores narrative nodes.
type NarrativeNodeDocPo struct {
	ID         string                 `bson:"id"`
	Title      string                 `bson:"title"`
	Body       string                 `bson:"body"`
	Options    []NarrativeOptionDocPo `bson:"options"`
	IsTerminal bool                   `bson:"is_terminal"`
}

// NarrativeOptionDocPo stores player options.
type NarrativeOptionDocPo struct {
	ID          string   `bson:"id"`
	Label       string   `bson:"label"`
	NextNodeID  string   `bson:"next_node_id"`
	EvidenceIDs []string `bson:"evidence_ids"`
	LeadsToEnd  bool     `bson:"leads_to_end"`
	SuccessEnd  bool     `bson:"success_end"`
}

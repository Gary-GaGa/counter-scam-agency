package intelligence

// EvidenceDocPo stores evidence data.
type EvidenceDocPo struct {
	ID              string `bson:"id"`
	Description     string `bson:"description"`
	Type            string `bson:"type"`
	IsContradiction bool   `bson:"is_contradiction"`
}

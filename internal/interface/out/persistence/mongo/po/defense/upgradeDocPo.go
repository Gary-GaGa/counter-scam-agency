package defense

// UpgradeDocPo represents a base upgrade MongoDB document.
type UpgradeDocPo struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Level       int    `bson:"level"`
	MaxLevel    int    `bson:"max_level"`
	Description string `bson:"description"`
}

package defense

// FacilityDocPo -
type FacilityDocPo struct {
	ID          string `bson:"id"`
	Type        string `bson:"type"`
	Name        string `bson:"name"`
	Level       int    `bson:"level"`
	MaxLevel    int    `bson:"max_level"`
	Description string `bson:"description"`
}

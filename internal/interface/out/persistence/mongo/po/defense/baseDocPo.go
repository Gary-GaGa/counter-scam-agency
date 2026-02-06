package defense

// BaseDocPo -
type BaseDocPo struct {
	ID            string          `bson:"id"`
	OwnerID       string          `bson:"owner_id"`
	SecurityLevel int             `bson:"security_level"`
	FacilitySlots int             `bson:"facility_slots"`
	Facilities    []FacilityDocPo `bson:"facilities"`
	Upgrades      []UpgradeDocPo  `bson:"upgrades"`
}

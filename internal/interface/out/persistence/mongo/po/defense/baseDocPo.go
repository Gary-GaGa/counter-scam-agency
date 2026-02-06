package defense

// BaseDocPo represents a defense base MongoDB document.
type BaseDocPo struct {
	ID            string          `bson:"id"`
	OwnerID       string          `bson:"owner_id"`
	SecurityLevel int             `bson:"security_level"`
	FacilitySlots int             `bson:"facility_slots"`
	Facilities    []FacilityDocPo `bson:"facilities"`
	Upgrades      []UpgradeDocPo  `bson:"upgrades"`
}

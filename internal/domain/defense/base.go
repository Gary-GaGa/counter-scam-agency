package defense

// FacilityType represents base facility category.
type FacilityType string

const (
	FacilityTypeFirewall FacilityType = "Firewall"
	FacilityTypeSIEM     FacilityType = "SIEM"
	FacilityTypeTraining FacilityType = "Training"
)

// Facility represents a buildable facility in the defense base.
type Facility struct {
	ID          string
	Type        FacilityType
	Name        string
	Level       int
	MaxLevel    int
	Description string
}

// Upgrade represents a base upgrade.
type Upgrade struct {
	ID          string
	Name        string
	Level       int
	MaxLevel    int
	Description string
}

// Base represents the digital defense base (Aggregate Root).
type Base struct {
	ID            string
	OwnerID       string
	SecurityLevel int
	FacilitySlots int
	Facilities    []Facility
	Upgrades      []Upgrade
}

// NewBase creates a base with default values.
func NewBase(id, ownerID string, slots int) *Base {
	if slots <= 0 {
		slots = 2
	}
	return &Base{
		ID:            id,
		OwnerID:       ownerID,
		SecurityLevel: 1,
		FacilitySlots: slots,
		Facilities:    make([]Facility, 0),
		Upgrades:      make([]Upgrade, 0),
	}
}

// AddFacility installs a new facility if there is slot capacity.
func (b *Base) AddFacility(f Facility) bool {
	if b == nil || f.ID == "" {
		return false
	}
	if len(b.Facilities) >= b.FacilitySlots {
		return false
	}
	for _, existing := range b.Facilities {
		if existing.ID == f.ID {
			return false
		}
	}
	if f.Level <= 0 {
		f.Level = 1
	}
	if f.MaxLevel <= 0 {
		f.MaxLevel = 3
	}
	b.Facilities = append(b.Facilities, f)
	return true
}

// UpgradeSecurity increases the base security level up to a cap.
func (b *Base) UpgradeSecurity(max int) bool {
	if b == nil {
		return false
	}
	if max <= 0 {
		max = 10
	}
	if b.SecurityLevel >= max {
		return false
	}
	b.SecurityLevel++
	return true
}

// UpgradeFacility increments facility level if possible.
func (b *Base) UpgradeFacility(facilityID string) bool {
	if b == nil {
		return false
	}
	for i := range b.Facilities {
		if b.Facilities[i].ID == facilityID {
			if b.Facilities[i].Level >= b.Facilities[i].MaxLevel {
				return false
			}
			b.Facilities[i].Level++
			return true
		}
	}
	return false
}

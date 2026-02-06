package dto

// BaseSummary represents a defense base for UI rendering.
type BaseSummary struct {
	ID            string
	OwnerID       string
	SecurityLevel int
	FacilitySlots int
	Facilities    []FacilitySummary
}

// FacilitySummary represents a facility for UI rendering.
type FacilitySummary struct {
	ID       string
	Type     string
	Name     string
	Level    int
	MaxLevel int
}

// FacilityInput represents a facility creation request.
type FacilityInput struct {
	ID          string
	Type        string
	Name        string
	Level       int
	MaxLevel    int
	Description string
}

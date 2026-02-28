package dto

// BaseSummary represents a defense base for UI rendering.
type BaseSummary struct {
	ID            string             `json:"id"`
	OwnerID       string             `json:"ownerId"`
	SecurityLevel int                `json:"securityLevel"`
	FacilitySlots int                `json:"facilitySlots"`
	Facilities    []FacilitySummary  `json:"facilities"`
}

// FacilitySummary represents a facility for UI rendering.
type FacilitySummary struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	MaxLevel int    `json:"maxLevel"`
}

// FacilityInput represents a facility creation request.
type FacilityInput struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	MaxLevel    int    `json:"maxLevel"`
	Description string `json:"description"`
}

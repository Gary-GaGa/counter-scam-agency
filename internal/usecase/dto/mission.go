package dto

// MissionSummary is used for mission listing screens.
type MissionSummary struct {
	ID               string
	Title            string
	Description      string
	Type             string
	Difficulty       int
	ReputationWeight int
}

// MissionDetail is used for mission detail and narrative rendering.
type MissionDetail struct {
	ID               string
	Title            string
	Description      string
	Type             string
	Difficulty       int
	ReputationWeight int
	Nodes            []NarrativeNode
	EvidenceList     []Evidence
}

// NarrativeNode represents a narrative node for UI rendering.
type NarrativeNode struct {
	ID         string
	Title      string
	Body       string
	Options    []NarrativeOption
	IsTerminal bool
}

// NarrativeOption represents a selectable option for UI rendering.
type NarrativeOption struct {
	ID          string
	Label       string
	NextNodeID  string
	EvidenceIDs []string
	LeadsToEnd  bool
	SuccessEnd  bool
}

// Evidence represents evidence summary for UI rendering.
type Evidence struct {
	ID              string
	Description     string
	Type            string
	IsContradiction bool
}

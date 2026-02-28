package dto

// MissionSummary is used for mission listing screens.
type MissionSummary struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	Difficulty       int    `json:"difficulty"`
	ReputationWeight int    `json:"reputationWeight"`
}

// MissionDetail is used for mission detail and narrative rendering.
type MissionDetail struct {
	ID               string           `json:"id"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	Type             string           `json:"type"`
	Difficulty       int              `json:"difficulty"`
	ReputationWeight int              `json:"reputationWeight"`
	VictimProfile    *VictimProfile   `json:"victimProfile,omitempty"`
	Nodes            []NarrativeNode  `json:"nodes"`
	EvidenceList     []Evidence       `json:"evidenceList"`
}

// VictimProfile represents victim psychology for UI rendering.
type VictimProfile struct {
	Anxiety   int    `json:"anxiety"`
	Trust     int    `json:"trust"`
	Urgency   int    `json:"urgency"`
	Isolation int    `json:"isolation"`
	RiskScore int    `json:"riskScore"`
	RiskLevel string `json:"riskLevel"`
}

// NarrativeNode represents a narrative node for UI rendering.
type NarrativeNode struct {
	ID         string             `json:"id"`
	Title      string             `json:"title"`
	Body       string             `json:"body"`
	Options    []NarrativeOption  `json:"options"`
	IsTerminal bool               `json:"isTerminal"`
}

// NarrativeOption represents a selectable option for UI rendering.
type NarrativeOption struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	NextNodeID  string   `json:"nextNodeId"`
	EvidenceIDs []string `json:"evidenceIds"`
	LeadsToEnd  bool     `json:"leadsToEnd"`
	SuccessEnd  bool     `json:"successEnd"`
}

// Evidence represents evidence summary for UI rendering.
type Evidence struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	Type            string `json:"type"`
	IsContradiction bool   `json:"isContradiction"`
}

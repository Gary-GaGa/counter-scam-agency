package dto

// InvestigationStartResult summarizes the start of an investigation.
type InvestigationStartResult struct {
	InvestigationID string `json:"investigationId"`
	PlayerID        string `json:"playerId"`
	MissionID       string `json:"missionId"`
	CurrentNodeID   string `json:"currentNodeId"`
	Status          string `json:"status"`
}

// NodeProgressResult summarizes the node transition result.
type NodeProgressResult struct {
	InvestigationID string `json:"investigationId"`
	MissionID       string `json:"missionId"`
	NodeID          string `json:"nodeId"`
	OptionID        string `json:"optionId"`
	NextNodeID      string `json:"nextNodeId"`
	Status          string `json:"status"`
}

// SubmitEvidenceResult summarizes the evidence submission outcome.
type SubmitEvidenceResult struct {
	InvestigationID  string `json:"investigationId"`
	EvidenceID       string `json:"evidenceId"`
	IsContradiction  bool   `json:"isContradiction"`
	AlreadyCollected bool   `json:"alreadyCollected"`
}

// CompleteResult summarizes the investigation completion outcome.
type CompleteResult struct {
	InvestigationID  string `json:"investigationId"`
	PlayerID         string `json:"playerId"`
	MissionID        string `json:"missionId"`
	Success          bool   `json:"success"`
	ReputationGained int    `json:"reputationGained"`
}

package dto

// InvestigationStartResult summarizes the start of an investigation.
type InvestigationStartResult struct {
	InvestigationID string
	PlayerID        string
	MissionID       string
	CurrentNodeID   string
	Status          string
}

// NodeProgressResult summarizes the node transition result.
type NodeProgressResult struct {
	InvestigationID string
	MissionID       string
	NodeID          string
	OptionID        string
	NextNodeID      string
	Status          string
}

// SubmitEvidenceResult summarizes the evidence submission outcome.
type SubmitEvidenceResult struct {
	InvestigationID  string
	EvidenceID       string
	IsContradiction  bool
	AlreadyCollected bool
}

// CompleteResult summarizes the investigation completion outcome.
type CompleteResult struct {
	InvestigationID  string
	PlayerID         string
	MissionID        string
	Success          bool
	ReputationGained int
}

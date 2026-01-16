package operation

// InvestigationStatus represents the current state of the investigation.
type InvestigationStatus string

const (
	InvestigationStatusActive    InvestigationStatus = "Active"
	InvestigationStatusCompleted InvestigationStatus = "Completed"
	InvestigationStatusFailed    InvestigationStatus = "Failed"
)

// Investigation represents the runtime state of a mission being played.
type Investigation struct {
	ID                   string
	PlayerID             string
	MissionID            string
	Status               InvestigationStatus
	CurrentNodeID        string
	NodeHistory          []NodeDecision
	CollectedEvidenceIDs []string
	SuspicionLevel       int // 0-100, represents how suspicious the scammer is of the player
}

// NewInvestigation starts a new investigation.
func NewInvestigation(id, playerID, missionID string) *Investigation {
	return &Investigation{
		ID:                   id,
		PlayerID:             playerID,
		MissionID:            missionID,
		Status:               InvestigationStatusActive,
		CurrentNodeID:        "",
		NodeHistory:          make([]NodeDecision, 0),
		CollectedEvidenceIDs: make([]string, 0),
		SuspicionLevel:       0,
	}
}

// NodeDecision records a player choice within a narrative node.
type NodeDecision struct {
	NodeID     string
	OptionID   string
	NextNodeID string
}

// SetCurrentNode updates the investigation's active node.
func (inv *Investigation) SetCurrentNode(nodeID string) {
	if inv.Status != InvestigationStatusActive {
		return
	}
	inv.CurrentNodeID = nodeID
}

// RecordDecision appends a node decision and advances the current node.
func (inv *Investigation) RecordDecision(nodeID, optionID, nextNodeID string) {
	if inv.Status != InvestigationStatusActive {
		return
	}
	inv.NodeHistory = append(inv.NodeHistory, NodeDecision{
		NodeID:     nodeID,
		OptionID:   optionID,
		NextNodeID: nextNodeID,
	})
	inv.CurrentNodeID = nextNodeID
}

// CollectEvidence records a piece of evidence as collected.
func (inv *Investigation) CollectEvidence(evidenceID string) {
	if inv.Status != InvestigationStatusActive {
		return
	}
	// Check for duplicates
	for _, id := range inv.CollectedEvidenceIDs {
		if id == evidenceID {
			return
		}
	}
	inv.CollectedEvidenceIDs = append(inv.CollectedEvidenceIDs, evidenceID)
}

// IncreaseSuspicion raises the suspicion level. If it hits 100, the investigation fails.
func (inv *Investigation) IncreaseSuspicion(amount int) {
	if inv.Status != InvestigationStatusActive {
		return
	}
	inv.SuspicionLevel += amount
	if inv.SuspicionLevel >= 100 {
		inv.SuspicionLevel = 100
		inv.Status = InvestigationStatusFailed
	}
}

// Complete marks the investigation as successfully solved.
func (inv *Investigation) Complete() {
	if inv.Status == InvestigationStatusActive {
		inv.Status = InvestigationStatusCompleted
	}
}

// Fail marks the investigation as failed.
func (inv *Investigation) Fail() {
	if inv.Status == InvestigationStatusActive {
		inv.Status = InvestigationStatusFailed
	}
}

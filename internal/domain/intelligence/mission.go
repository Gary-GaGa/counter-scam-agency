package intelligence

// Mission represents a scam investigation scenario and acts as the Aggregate Root for the Intelligence Context.
type Mission struct {
	ID               string
	Title            string
	Description      string
	Type             ScamType
	Difficulty       int
	ReputationWeight int
	Nodes            []NarrativeNode
	EvidenceList     []Evidence
}

// NewMission creates a new mission.
func NewMission(id, title, description string, scamType ScamType, difficulty int, reputationWeight int) *Mission {
	return &Mission{
		ID:               id,
		Title:            title,
		Description:      description,
		Type:             scamType,
		Difficulty:       difficulty,
		ReputationWeight: reputationWeight,
		Nodes:            make([]NarrativeNode, 0),
		EvidenceList:     make([]Evidence, 0),
	}
}

// NarrativeNode represents a node in the text adventure flow.
type NarrativeNode struct {
	ID         string
	Title      string
	Body       string
	Options    []NarrativeOption
	IsTerminal bool
}

// NarrativeOption represents a player choice from a node.
type NarrativeOption struct {
	ID          string
	Label       string
	NextNodeID  string
	EvidenceIDs []string
	LeadsToEnd  bool
	SuccessEnd  bool
}

// AddEvidence appends a piece of evidence to the mission.
func (m *Mission) AddEvidence(ev Evidence) {
	m.EvidenceList = append(m.EvidenceList, ev)
}

// AddNode appends a narrative node to the mission.
func (m *Mission) AddNode(node NarrativeNode) {
	m.Nodes = append(m.Nodes, node)
}

// GetNode retrieves a narrative node by ID.
func (m *Mission) GetNode(id string) *NarrativeNode {
	for i := range m.Nodes {
		if m.Nodes[i].ID == id {
			return &m.Nodes[i]
		}
	}
	return nil
}

// GetEvidence retrieves a piece of evidence by ID.
func (m *Mission) GetEvidence(id string) *Evidence {
	for i := range m.EvidenceList {
		if m.EvidenceList[i].ID == id {
			return &m.EvidenceList[i]
		}
	}
	return nil
}

// ValidateContradiction checks if the provided evidence ID corresponds to a contradiction.
func (m *Mission) ValidateContradiction(evidenceID string) bool {
	ev := m.GetEvidence(evidenceID)
	if ev != nil {
		return ev.IsContradiction
	}
	return false
}

// ReputationGain returns the reputation gained for this mission.
// Failed missions return 0 reputation.
func (m *Mission) ReputationGain(success bool) int {
	if !success {
		return 0
	}
	weight := m.ReputationWeight
	if weight <= 0 {
		return 0
	}
	if m.Difficulty <= 0 {
		return weight
	}
	return weight * m.Difficulty
}

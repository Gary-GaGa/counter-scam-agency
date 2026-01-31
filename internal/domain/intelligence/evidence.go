package intelligence

// EvidenceType classifies the nature of the evidence.
type EvidenceType string

const (
	EvidenceTypeDialogue    EvidenceType = "Dialogue"
	EvidenceTypeTransaction EvidenceType = "Transaction"
	EvidenceTypeImage       EvidenceType = "Image"
	EvidenceTypeDocument    EvidenceType = "Document"
)

// Evidence represents a clue or proof collected during investigation.
type Evidence struct {
	ID              string
	Description     string
	Type            EvidenceType
	IsContradiction bool // True if this piece of evidence directly contradicts a lie
}

// NewEvidence creates a new piece of evidence.
func NewEvidence(id, description string, evType EvidenceType, isContradiction bool) *Evidence {
	return &Evidence{
		ID:              id,
		Description:     description,
		Type:            evType,
		IsContradiction: isContradiction,
	}
}

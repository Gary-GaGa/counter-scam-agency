package personnel

// SkillType represents an active AI skill category.
type SkillType string

const (
	SkillTypeAnalysis  SkillType = "Analysis"
	SkillTypeNegotiation SkillType = "Negotiation"
	SkillTypeDefense   SkillType = "Defense"
	SkillTypeForensics SkillType = "Forensics"
)

// Skill represents an active ability for the AI partner.
type Skill struct {
	ID                 string
	Type               SkillType
	Name               string
	Description        string
	CooldownSeconds    int
	ReputationRequired int
	RequiredModuleIDs  []string
}

// RequiresModules indicates whether the skill depends on modules.
func (s Skill) RequiresModules() bool {
	return len(s.RequiredModuleIDs) > 0
}

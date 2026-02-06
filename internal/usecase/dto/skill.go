package dto

// SkillSummary represents an AI skill for UI rendering.
type SkillSummary struct {
	ID                string
	Type              string
	Name              string
	Description       string
	CooldownSeconds   int
	RequiredModuleIDs []string
	Unlocked          bool
	Equipped          bool
	CooldownRemaining int
}

// SkillActionResult summarizes a skill action outcome.
type SkillActionResult struct {
	PlayerID          string
	SkillID           string
	Unlocked          bool
	Equipped          bool
	CooldownRemaining int
}

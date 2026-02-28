package dto

// PlayerSummary represents a player for UI rendering.
type PlayerSummary struct {
	ID         string     `json:"id"`
	Reputation int        `json:"reputation"`
	Stats      StatsSummary `json:"stats"`
}

// StatsSummary represents player stats for UI rendering.
type StatsSummary struct {
	Logic      int `json:"logic"`
	Tech       int `json:"tech"`
	Charisma   int `json:"charisma"`
	Resilience int `json:"resilience"`
}

// SkillSummary represents an AI skill for UI rendering.
type SkillSummary struct {
	ID                string   `json:"id"`
	Type              string   `json:"type"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	CooldownSeconds   int      `json:"cooldownSeconds"`
	RequiredModuleIDs []string `json:"requiredModuleIds"`
	Unlocked          bool     `json:"unlocked"`
	Equipped          bool     `json:"equipped"`
	CooldownRemaining int      `json:"cooldownRemaining"`
}

// SkillActionResult summarizes a skill action outcome.
type SkillActionResult struct {
	PlayerID          string `json:"playerId"`
	SkillID           string `json:"skillId"`
	Unlocked          bool   `json:"unlocked"`
	Equipped          bool   `json:"equipped"`
	CooldownRemaining int    `json:"cooldownRemaining"`
}

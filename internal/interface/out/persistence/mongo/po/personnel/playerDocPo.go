package personnel

// PlayerDocPo stores player data in MongoDB.
type PlayerDocPo struct {
	ID                string          `bson:"id"`
	Stats             StatsDocPo      `bson:"stats"`
	Partner           *AIPartnerDocPo `bson:"partner"`
	Reputation        int             `bson:"reputation"`
	UnlockedModuleIDs []string        `bson:"unlocked_module_ids"`
	UnlockedSkillIDs  []string        `bson:"unlocked_skill_ids"`
}

// AIPartnerDocPo stores AI partner data.
type AIPartnerDocPo struct {
	Loadout        []ModuleDocPo  `bson:"loadout"`
	Skills         []SkillDocPo   `bson:"skills"`
	SkillCooldowns map[string]int `bson:"skill_cooldowns"`
	Personality    string         `bson:"personality"`
}

// ModuleDocPo stores module data.
type ModuleDocPo struct {
	ID                 string     `bson:"id"`
	Type               string     `bson:"type"`
	Name               string     `bson:"name"`
	Description        string     `bson:"description"`
	Level              int        `bson:"level"`
	ReputationRequired int        `bson:"reputation_required"`
	StatBonus          StatsDocPo `bson:"stat_bonus"`
}

// SkillDocPo stores skill data.
type SkillDocPo struct {
	ID                 string   `bson:"id"`
	Type               string   `bson:"type"`
	Name               string   `bson:"name"`
	Description        string   `bson:"description"`
	CooldownSeconds    int      `bson:"cooldown_seconds"`
	ReputationRequired int      `bson:"reputation_required"`
	RequiredModuleIDs  []string `bson:"required_module_ids"`
}

// StatsDocPo stores stats data.
type StatsDocPo struct {
	Logic      int `bson:"logic"`
	Tech       int `bson:"tech"`
	Charisma   int `bson:"charisma"`
	Resilience int `bson:"resilience"`
}

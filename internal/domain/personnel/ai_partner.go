package personnel

// ModuleType represents the type of module installed on the AI.
type ModuleType string

const (
	ModuleVoiceAnalyzer  ModuleType = "VoiceAnalyzer"
	ModuleCryptoTracer   ModuleType = "CryptoTracer"
	ModuleEmpathyEngine  ModuleType = "EmpathyEngine"
	ModuleMentalFirewall ModuleType = "MentalFirewall"
)

// Module represents an upgradeable component for the AI.
type Module struct {
	ID                 string
	Type               ModuleType
	Name               string
	Description        string
	Level              int
	ReputationRequired int
	StatBonus          Stats
}

// PersonalityType defines the AI's current behavioral mode.
type PersonalityType string

const (
	PersonalityLogicDriven   PersonalityType = "LogicDriven"
	PersonalityEmpathyDriven PersonalityType = "EmpathyDriven"
	PersonalityBalanced      PersonalityType = "Balanced"
)

// AIPartner represents the player's AI companion.
type AIPartner struct {
	Loadout        []Module
	Skills         []Skill
	SkillCooldowns map[string]int
	Personality    PersonalityType
}

// NewAIPartner creates a default AI partner.
func NewAIPartner() *AIPartner {
	return &AIPartner{
		Loadout:        make([]Module, 0),
		Skills:         make([]Skill, 0),
		SkillCooldowns: make(map[string]int),
		Personality:    PersonalityBalanced,
	}
}

// InstallModule adds a module to the AI's loadout.
func (ai *AIPartner) InstallModule(mod Module) {
	ai.Loadout = append(ai.Loadout, mod)
}

// LearnSkill adds a skill to the AI loadout if not already learned.
func (ai *AIPartner) LearnSkill(skill Skill) bool {
	if ai.HasSkill(skill.ID) {
		return false
	}
	ai.Skills = append(ai.Skills, skill)
	if ai.SkillCooldowns == nil {
		ai.SkillCooldowns = make(map[string]int)
	}
	if _, ok := ai.SkillCooldowns[skill.ID]; !ok {
		ai.SkillCooldowns[skill.ID] = 0
	}
	return true
}

// HasSkill checks if the AI has learned a skill.
func (ai *AIPartner) HasSkill(skillID string) bool {
	for _, skill := range ai.Skills {
		if skill.ID == skillID {
			return true
		}
	}
	return false
}

// HasModules checks whether all required modules are installed.
func (ai *AIPartner) HasModules(requiredIDs []string) bool {
	if len(requiredIDs) == 0 {
		return true
	}
	installed := make(map[string]struct{}, len(ai.Loadout))
	for _, mod := range ai.Loadout {
		installed[mod.ID] = struct{}{}
	}
	for _, id := range requiredIDs {
		if _, ok := installed[id]; !ok {
			return false
		}
	}
	return true
}

// CooldownRemaining returns remaining cooldown seconds.
func (ai *AIPartner) CooldownRemaining(skillID string) int {
	if ai.SkillCooldowns == nil {
		return 0
	}
	return ai.SkillCooldowns[skillID]
}

// UseSkill activates a skill and sets cooldown if possible.
func (ai *AIPartner) UseSkill(skillID string) bool {
	skill := ai.GetSkill(skillID)
	if skill == nil {
		return false
	}
	if !ai.HasModules(skill.RequiredModuleIDs) {
		return false
	}
	if ai.CooldownRemaining(skillID) > 0 {
		return false
	}
	if ai.SkillCooldowns == nil {
		ai.SkillCooldowns = make(map[string]int)
	}
	if skill.CooldownSeconds < 0 {
		ai.SkillCooldowns[skillID] = 0
		return true
	}
	ai.SkillCooldowns[skillID] = skill.CooldownSeconds
	return true
}

// TickCooldown reduces all cooldowns by seconds.
func (ai *AIPartner) TickCooldown(seconds int) {
	if seconds <= 0 || ai.SkillCooldowns == nil {
		return
	}
	for id, remaining := range ai.SkillCooldowns {
		next := remaining - seconds
		if next < 0 {
			next = 0
		}
		ai.SkillCooldowns[id] = next
	}
}

// GetSkill retrieves a learned skill by ID.
func (ai *AIPartner) GetSkill(skillID string) *Skill {
	for i := range ai.Skills {
		if ai.Skills[i].ID == skillID {
			return &ai.Skills[i]
		}
	}
	return nil
}

// GetTotalBonus calculates the total stats bonus from all installed modules.
func (ai *AIPartner) GetTotalBonus() Stats {
	total := Stats{}
	for _, mod := range ai.Loadout {
		total = total.Add(mod.StatBonus)
	}
	return total
}

// InstalledModuleIDs returns the IDs of all installed modules.
func (ai *AIPartner) InstalledModuleIDs() []string {
	ids := make([]string, len(ai.Loadout))
	for i, mod := range ai.Loadout {
		ids[i] = mod.ID
	}
	return ids
}

// LearnedSkillIDs returns the IDs of all learned skills.
func (ai *AIPartner) LearnedSkillIDs() []string {
	ids := make([]string, len(ai.Skills))
	for i, skill := range ai.Skills {
		ids[i] = skill.ID
	}
	return ids
}

package personnel

// Player represents the main character and acts as the Aggregate Root for the Personnel Context.
type Player struct {
	ID                string
	Stats             Stats
	Partner           *AIPartner
	Reputation        int
	UnlockedModuleIDs map[string]struct{}
	UnlockedSkillIDs  map[string]struct{}
}

// NewPlayer creates a player with initial stats.
func NewPlayer(id string) *Player {
	return &Player{
		ID: id,
		Stats: Stats{
			Logic:      10,
			Tech:       10,
			Charisma:   10,
			Resilience: 10,
		},
		Partner:           NewAIPartner(),
		Reputation:        0,
		UnlockedModuleIDs: make(map[string]struct{}),
		UnlockedSkillIDs:  make(map[string]struct{}),
	}
}

// GetTotalStats returns the player's base stats plus any bonuses from the AI Partner.
func (p *Player) GetTotalStats() Stats {
	if p.Partner == nil {
		return p.Stats
	}
	return p.Stats.Add(p.Partner.GetTotalBonus())
}

// AddReputation increases the player's reputation.
func (p *Player) AddReputation(amount int) {
	if amount <= 0 {
		return
	}
	p.Reputation += amount
}

// AddStats increases the player's base stats. Negative values are ignored per field.
func (p *Player) AddStats(logic, tech, charisma, resilience int) {
	if logic > 0 {
		p.Stats.Logic += logic
	}
	if tech > 0 {
		p.Stats.Tech += tech
	}
	if charisma > 0 {
		p.Stats.Charisma += charisma
	}
	if resilience > 0 {
		p.Stats.Resilience += resilience
	}
}

// IsModuleUnlocked checks whether a module has been unlocked.
func (p *Player) IsModuleUnlocked(moduleID string) bool {
	if len(p.UnlockedModuleIDs) == 0 {
		return false
	}
	_, ok := p.UnlockedModuleIDs[moduleID]
	return ok
}

// CanUnlockModule checks if the player meets the reputation requirement.
func (p *Player) CanUnlockModule(mod Module) bool {
	if mod.ID == "" {
		return false
	}
	if p.IsModuleUnlocked(mod.ID) {
		return false
	}
	return p.Reputation >= mod.ReputationRequired
}

// UnlockModule unlocks a module if requirements are met.
func (p *Player) UnlockModule(mod Module) bool {
	if !p.CanUnlockModule(mod) {
		return false
	}
	if p.UnlockedModuleIDs == nil {
		p.UnlockedModuleIDs = make(map[string]struct{})
	}
	p.UnlockedModuleIDs[mod.ID] = struct{}{}
	return true
}

// EquipPartnerModule installs a module onto the AI Partner if unlocked.
func (p *Player) EquipPartnerModule(mod Module) bool {
	if p.Partner == nil {
		return false
	}
	if !p.IsModuleUnlocked(mod.ID) {
		return false
	}
	p.Partner.InstallModule(mod)
	return true
}

// IsSkillUnlocked checks whether a skill has been unlocked.
func (p *Player) IsSkillUnlocked(skillID string) bool {
	if len(p.UnlockedSkillIDs) == 0 {
		return false
	}
	_, ok := p.UnlockedSkillIDs[skillID]
	return ok
}

// CanUnlockSkill checks if the player meets the reputation requirement.
func (p *Player) CanUnlockSkill(skill Skill) bool {
	if skill.ID == "" {
		return false
	}
	if p.IsSkillUnlocked(skill.ID) {
		return false
	}
	return p.Reputation >= skill.ReputationRequired
}

// UnlockSkill unlocks a skill if requirements are met.
func (p *Player) UnlockSkill(skill Skill) bool {
	if !p.CanUnlockSkill(skill) {
		return false
	}
	if p.UnlockedSkillIDs == nil {
		p.UnlockedSkillIDs = make(map[string]struct{})
	}
	p.UnlockedSkillIDs[skill.ID] = struct{}{}
	return true
}

// EquipPartnerSkill installs a skill onto the AI Partner if unlocked.
func (p *Player) EquipPartnerSkill(skill Skill) bool {
	if p.Partner == nil {
		return false
	}
	if !p.IsSkillUnlocked(skill.ID) {
		return false
	}
	return p.Partner.LearnSkill(skill)
}

// ActivatePartnerSkill activates a learned skill on the AI partner.
func (p *Player) ActivatePartnerSkill(skillID string) bool {
	if p.Partner == nil {
		return false
	}
	return p.Partner.UseSkill(skillID)
}

// TickPartnerSkillCooldowns reduces cooldowns for all AI skills.
func (p *Player) TickPartnerSkillCooldowns(seconds int) {
	if p.Partner == nil {
		return
	}
	p.Partner.TickCooldown(seconds)
}

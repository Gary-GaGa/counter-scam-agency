package personnel

// Player represents the main character and acts as the Aggregate Root for the Personnel Context.
type Player struct {
	ID                string
	Stats             Stats
	Partner           *AIPartner
	Reputation        int
	UnlockedModuleIDs map[string]struct{}
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

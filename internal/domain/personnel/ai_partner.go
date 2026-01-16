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
	Loadout     []Module
	Personality PersonalityType
}

// NewAIPartner creates a default AI partner.
func NewAIPartner() *AIPartner {
	return &AIPartner{
		Loadout:     make([]Module, 0),
		Personality: PersonalityBalanced,
	}
}

// InstallModule adds a module to the AI's loadout.
func (ai *AIPartner) InstallModule(mod Module) {
	ai.Loadout = append(ai.Loadout, mod)
}

// GetTotalBonus calculates the total stats bonus from all installed modules.
func (ai *AIPartner) GetTotalBonus() Stats {
	total := Stats{}
	for _, mod := range ai.Loadout {
		total = total.Add(mod.StatBonus)
	}
	return total
}

package personnel_test

import (
	"testing"

	"counter-scam-agency/internal/domain/personnel"

	"github.com/stretchr/testify/assert"
)

func TestPlayerStatsCalculation(t *testing.T) {
	// 1. Create Player
	player := personnel.NewPlayer("agent-007")

	// 2. Verify Initial Stats
	initialStats := player.GetTotalStats()
	assert.Equal(t, 10, initialStats.Logic)
	assert.Equal(t, 10, initialStats.Tech)

	// 3. Create a Module with Bonus
	mod := personnel.Module{
		ID:                 "mod-001",
		Type:               personnel.ModuleVoiceAnalyzer,
		Name:               "Basic Voice Analyzer",
		ReputationRequired: 50,
		StatBonus: personnel.Stats{
			Logic: 5,
		},
	}

	// 4. Unlock and Equip Module
	player.AddReputation(100)
	assert.True(t, player.UnlockModule(mod))
	assert.True(t, player.EquipPartnerModule(mod))

	// 5. Verify Updated Stats (Base + Bonus)
	newStats := player.GetTotalStats()
	assert.Equal(t, 15, newStats.Logic) // 10 + 5
	assert.Equal(t, 10, newStats.Tech)  // Unchanged
}

func TestAIPartnerPersonality(t *testing.T) {
	ai := personnel.NewAIPartner()
	assert.Equal(t, personnel.PersonalityBalanced, ai.Personality)
}

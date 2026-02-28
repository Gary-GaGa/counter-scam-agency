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

func TestPlayerAddStats(t *testing.T) {
	player := personnel.NewPlayer("agent-009")
	assert.Equal(t, 10, player.Stats.Logic)
	assert.Equal(t, 10, player.Stats.Tech)

	player.AddStats(3, 0, 2, 1)
	assert.Equal(t, 13, player.Stats.Logic)
	assert.Equal(t, 10, player.Stats.Tech)
	assert.Equal(t, 12, player.Stats.Charisma)
	assert.Equal(t, 11, player.Stats.Resilience)

	// 負數不影響
	player.AddStats(-5, -1, -1, -1)
	assert.Equal(t, 13, player.Stats.Logic)
}

func TestPlayerSkillUnlockAndEquip(t *testing.T) {
	player := personnel.NewPlayer("agent-008")
	skill := personnel.Skill{
		ID:                 "skill-001",
		Type:               personnel.SkillTypeAnalysis,
		Name:               "Rapid Analysis",
		CooldownSeconds:    10,
		ReputationRequired: 30,
	}

	player.AddReputation(30)
	assert.True(t, player.UnlockSkill(skill))
	assert.True(t, player.EquipPartnerSkill(skill))
	assert.True(t, player.Partner.HasSkill("skill-001"))
	assert.True(t, player.ActivatePartnerSkill("skill-001"))
	assert.False(t, player.ActivatePartnerSkill("skill-001"))
	player.TickPartnerSkillCooldowns(10)
	assert.True(t, player.ActivatePartnerSkill("skill-001"))
}

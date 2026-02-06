package intelligence_test

import (
	"testing"

	"counter-scam-agency/internal/domain/intelligence"

	"github.com/stretchr/testify/assert"
)

func TestMissionInitialization(t *testing.T) {
	mission := intelligence.NewMission(
		"mission-001",
		"The Nigerian Prince",
		"Classic email scam",
		intelligence.ScamTypePhishing,
		2,
		0,
	)

	assert.NotNil(t, mission)
	assert.Equal(t, "mission-001", mission.ID)
	assert.Equal(t, intelligence.ScamTypePhishing, mission.Type)
	assert.Empty(t, mission.EvidenceList)
}

func TestEvidenceManagement(t *testing.T) {
	mission := intelligence.NewMission("m-001", "Test", "Desc", intelligence.ScamTypeInvestment, 1, 0)

	// 1. Create Evidence
	ev1 := intelligence.NewEvidence("ev-001", "Fake Bank Statement", intelligence.EvidenceTypeDocument, true)
	ev2 := intelligence.NewEvidence("ev-002", "Normal Chat Log", intelligence.EvidenceTypeDialogue, false)

	// 2. Add Evidence
	mission.AddEvidence(*ev1)
	mission.AddEvidence(*ev2)

	// 3. Verify Retrieval
	retrievedEv := mission.GetEvidence("ev-001")
	assert.NotNil(t, retrievedEv)
	assert.Equal(t, "Fake Bank Statement", retrievedEv.Description)

	// 4. Verify Contradiction Validation
	assert.True(t, mission.ValidateContradiction("ev-001"))
	assert.False(t, mission.ValidateContradiction("ev-002"))
	assert.False(t, mission.ValidateContradiction("non-existent-id"))
}

func TestVictimProfileRiskLevel(t *testing.T) {
	profile := intelligence.VictimProfile{
		Anxiety:   80,
		Trust:     70,
		Urgency:   60,
		Isolation: 50,
	}
	assert.Equal(t, 65, profile.RiskScore())
	assert.Equal(t, intelligence.RiskLevelMedium, profile.RiskLevel())
}

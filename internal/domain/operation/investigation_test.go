package operation_test

import (
	"testing"

	"counter-scam-agency/internal/domain/operation"

	"github.com/stretchr/testify/assert"
)

func TestInvestigationLifecycle(t *testing.T) {
	inv := operation.NewInvestigation("inv-001", "player-001", "mission-001")

	// 1. Initial State
	assert.Equal(t, operation.InvestigationStatusActive, inv.Status)
	assert.Equal(t, 0, inv.SuspicionLevel)
	assert.Empty(t, inv.CollectedEvidenceIDs)

	// 2. Collect Evidence
	inv.CollectEvidence("ev-001")
	assert.Contains(t, inv.CollectedEvidenceIDs, "ev-001")

	// 3. Duplicate Evidence Check
	inv.CollectEvidence("ev-001")
	assert.Len(t, inv.CollectedEvidenceIDs, 1)

	// 4. Complete Investigation
	inv.Complete()
	assert.Equal(t, operation.InvestigationStatusCompleted, inv.Status)

	// 5. Try to modify after completion (should fail)
	inv.CollectEvidence("ev-002")
	assert.NotContains(t, inv.CollectedEvidenceIDs, "ev-002")
}

func TestSuspicionMechanic(t *testing.T) {
	inv := operation.NewInvestigation("inv-002", "p-1", "m-1")

	// 1. Increase Suspicion
	inv.IncreaseSuspicion(50)
	assert.Equal(t, 50, inv.SuspicionLevel)
	assert.Equal(t, operation.InvestigationStatusActive, inv.Status)

	// 2. Increase to Failure
	inv.IncreaseSuspicion(60)                // Total 110
	assert.Equal(t, 100, inv.SuspicionLevel) // Capped at 100
	assert.Equal(t, operation.InvestigationStatusFailed, inv.Status)
}

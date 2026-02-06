package defense_test

import (
	"testing"

	"counter-scam-agency/internal/domain/defense"

	"github.com/stretchr/testify/assert"
)

func TestBaseLifecycle(t *testing.T) {
	base := defense.NewBase("base-1", "player-1", 2)
	assert.Equal(t, 1, base.SecurityLevel)

	facility := defense.Facility{ID: "fac-1", Type: defense.FacilityTypeFirewall, Name: "Starter Firewall"}
	assert.True(t, base.AddFacility(facility))
	assert.False(t, base.AddFacility(facility))

	assert.True(t, base.UpgradeFacility("fac-1"))
	assert.True(t, base.UpgradeSecurity(2))
	assert.False(t, base.UpgradeSecurity(2))
}

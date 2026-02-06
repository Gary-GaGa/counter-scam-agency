package defense_test

import (
	"context"
	"testing"

	domainDefense "counter-scam-agency/internal/domain/defense"
	usecaseDefense "counter-scam-agency/internal/usecase/defense"
	"counter-scam-agency/internal/usecase/dto"

	"github.com/stretchr/testify/assert"
)

type baseRepo struct {
	bases map[string]*domainDefense.Base
}

func (r *baseRepo) Save(_ context.Context, base *domainDefense.Base) error {
	r.bases[base.ID] = base
	return nil
}

func (r *baseRepo) FindByID(_ context.Context, id string) (*domainDefense.Base, error) {
	return r.bases[id], nil
}

func (r *baseRepo) FindByOwnerID(_ context.Context, ownerID string) ([]*domainDefense.Base, error) {
	results := make([]*domainDefense.Base, 0)
	for _, base := range r.bases {
		if base.OwnerID == ownerID {
			results = append(results, base)
		}
	}
	return results, nil
}

func TestCreateAndUpgradeBase(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	created, err := svc.CreateBase(context.Background(), "base-1", "player-1", 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.SecurityLevel)

	updated, err := svc.UpgradeSecurity(context.Background(), "base-1", 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, updated.SecurityLevel)
}

func TestAddFacility(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.CreateBase(context.Background(), "base-1", "player-1", 1)
	assert.NoError(t, err)

	updated, err := svc.AddFacility(context.Background(), "base-1", dto.FacilityInput{
		ID:       "fac-1",
		Type:     "Firewall",
		Name:     "Starter Firewall",
		Level:    1,
		MaxLevel: 2,
	})
	assert.NoError(t, err)
	assert.Len(t, updated.Facilities, 1)
}

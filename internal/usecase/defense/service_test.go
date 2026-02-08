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

func TestGetBase(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.CreateBase(context.Background(), "base-1", "player-1", 2)
	assert.NoError(t, err)

	result, err := svc.GetBase(context.Background(), "base-1")
	assert.NoError(t, err)
	assert.Equal(t, "base-1", result.ID)
	assert.Equal(t, "player-1", result.OwnerID)
}

func TestGetBase_NotFound(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.GetBase(context.Background(), "missing")
	assert.ErrorIs(t, err, usecaseDefense.ErrBaseNotFound)
}

func TestCreateBase_Duplicate(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.CreateBase(context.Background(), "base-1", "player-1", 2)
	assert.NoError(t, err)

	_, err = svc.CreateBase(context.Background(), "base-1", "player-1", 2)
	assert.ErrorIs(t, err, usecaseDefense.ErrBaseExists)
}

func TestUpgradeFacility(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.CreateBase(context.Background(), "base-1", "player-1", 2)
	assert.NoError(t, err)

	_, err = svc.AddFacility(context.Background(), "base-1", dto.FacilityInput{
		ID:       "fac-1",
		Type:     "Firewall",
		Name:     "Starter Firewall",
		Level:    1,
		MaxLevel: 3,
	})
	assert.NoError(t, err)

	result, err := svc.UpgradeFacility(context.Background(), "base-1", "fac-1")
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Facilities[0].Level)
}

func TestUpgradeFacility_BaseNotFound(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.UpgradeFacility(context.Background(), "missing", "fac-1")
	assert.ErrorIs(t, err, usecaseDefense.ErrBaseNotFound)
}

func TestAddFacility_Full(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.CreateBase(context.Background(), "base-1", "player-1", 1)
	assert.NoError(t, err)

	_, err = svc.AddFacility(context.Background(), "base-1", dto.FacilityInput{
		ID: "fac-1", Type: "Firewall", Name: "FW", Level: 1, MaxLevel: 2,
	})
	assert.NoError(t, err)

	_, err = svc.AddFacility(context.Background(), "base-1", dto.FacilityInput{
		ID: "fac-2", Type: "Scanner", Name: "SC", Level: 1, MaxLevel: 2,
	})
	assert.ErrorIs(t, err, usecaseDefense.ErrFacilityFull)
}

func TestUpgradeSecurity_BaseNotFound(t *testing.T) {
	repo := &baseRepo{bases: map[string]*domainDefense.Base{}}
	svc := usecaseDefense.NewService(repo)

	_, err := svc.UpgradeSecurity(context.Background(), "missing", 5)
	assert.ErrorIs(t, err, usecaseDefense.ErrBaseNotFound)
}

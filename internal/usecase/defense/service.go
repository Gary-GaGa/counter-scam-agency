package defense

import (
	"context"
	"errors"
	"fmt"

	"counter-scam-agency/internal/domain/defense"
	"counter-scam-agency/internal/usecase/dto"
	portin "counter-scam-agency/internal/usecase/port/in"
)

// Compile-time interface guard.
var _ portin.DefenseUsecase = (*Service)(nil)

var (
	ErrBaseNotFound = errors.New("base not found")
	ErrBaseExists   = errors.New("base already exists")
	ErrFacilityFull = errors.New("facility slots full")
	ErrFacilityFail = errors.New("facility operation failed")
)

// Service orchestrates defense base flows.
type Service struct {
	bases defense.BaseRepository
}

// NewService creates a new defense service.
func NewService(bases defense.BaseRepository) *Service {
	return &Service{bases: bases}
}

// CreateBase creates a new base for a player.
func (s *Service) CreateBase(ctx context.Context, baseID, ownerID string, slots int) (*dto.BaseSummary, error) {
	if baseID != "" {
		existing, err := s.bases.FindByID(ctx, baseID)
		if err != nil {
			return nil, fmt.Errorf("find base: %w", err)
		}
		if existing != nil {
			return nil, ErrBaseExists
		}
	}
	base := defense.NewBase(baseID, ownerID, slots)
	if err := s.bases.Save(ctx, base); err != nil {
		return nil, fmt.Errorf("save base: %w", err)
	}
	return mapBase(base), nil
}

// GetBase retrieves a base.
func (s *Service) GetBase(ctx context.Context, baseID string) (*dto.BaseSummary, error) {
	base, err := s.bases.FindByID(ctx, baseID)
	if err != nil {
		return nil, fmt.Errorf("find base: %w", err)
	}
	if base == nil {
		return nil, ErrBaseNotFound
	}
	return mapBase(base), nil
}

// AddFacility adds a facility to the base.
func (s *Service) AddFacility(ctx context.Context, baseID string, facility dto.FacilityInput) (*dto.BaseSummary, error) {
	base, err := s.bases.FindByID(ctx, baseID)
	if err != nil {
		return nil, fmt.Errorf("find base: %w", err)
	}
	if base == nil {
		return nil, ErrBaseNotFound
	}

	added := base.AddFacility(defense.Facility{
		ID:          facility.ID,
		Type:        defense.FacilityType(facility.Type),
		Name:        facility.Name,
		Level:       facility.Level,
		MaxLevel:    facility.MaxLevel,
		Description: facility.Description,
	})
	if !added {
		return nil, ErrFacilityFull
	}
	if err := s.bases.Save(ctx, base); err != nil {
		return nil, fmt.Errorf("save base: %w", err)
	}
	return mapBase(base), nil
}

// UpgradeSecurity increases the base security level.
func (s *Service) UpgradeSecurity(ctx context.Context, baseID string, maxLevel int) (*dto.BaseSummary, error) {
	base, err := s.bases.FindByID(ctx, baseID)
	if err != nil {
		return nil, fmt.Errorf("find base: %w", err)
	}
	if base == nil {
		return nil, ErrBaseNotFound
	}
	if !base.UpgradeSecurity(maxLevel) {
		return nil, ErrFacilityFail
	}
	if err := s.bases.Save(ctx, base); err != nil {
		return nil, fmt.Errorf("save base: %w", err)
	}
	return mapBase(base), nil
}

// UpgradeFacility upgrades a facility by ID.
func (s *Service) UpgradeFacility(ctx context.Context, baseID, facilityID string) (*dto.BaseSummary, error) {
	base, err := s.bases.FindByID(ctx, baseID)
	if err != nil {
		return nil, fmt.Errorf("find base: %w", err)
	}
	if base == nil {
		return nil, ErrBaseNotFound
	}
	if !base.UpgradeFacility(facilityID) {
		return nil, ErrFacilityFail
	}
	if err := s.bases.Save(ctx, base); err != nil {
		return nil, fmt.Errorf("save base: %w", err)
	}
	return mapBase(base), nil
}

func mapBase(base *defense.Base) *dto.BaseSummary {
	if base == nil {
		return nil
	}
	facilities := make([]dto.FacilitySummary, 0, len(base.Facilities))
	for _, facility := range base.Facilities {
		facilities = append(facilities, dto.FacilitySummary{
			ID:       facility.ID,
			Type:     string(facility.Type),
			Name:     facility.Name,
			Level:    facility.Level,
			MaxLevel: facility.MaxLevel,
		})
	}
	return &dto.BaseSummary{
		ID:            base.ID,
		OwnerID:       base.OwnerID,
		SecurityLevel: base.SecurityLevel,
		FacilitySlots: base.FacilitySlots,
		Facilities:    facilities,
	}
}

package in

import (
	"context"

	"counter-scam-agency/internal/usecase/dto"
)

// DefenseUsecase defines input port for defense base flows.
type DefenseUsecase interface {
	CreateBase(ctx context.Context, baseID, ownerID string, slots int) (*dto.BaseSummary, error)
	GetBase(ctx context.Context, baseID string) (*dto.BaseSummary, error)
	AddFacility(ctx context.Context, baseID string, facility dto.FacilityInput) (*dto.BaseSummary, error)
	UpgradeSecurity(ctx context.Context, baseID string, maxLevel int) (*dto.BaseSummary, error)
	UpgradeFacility(ctx context.Context, baseID, facilityID string) (*dto.BaseSummary, error)
}

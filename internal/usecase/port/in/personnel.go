package in

import (
	"context"

	"counter-scam-agency/internal/usecase/dto"
)

// PersonnelUsecase defines input port for player/AI skill flows.
type PersonnelUsecase interface {
	CreatePlayer(ctx context.Context, playerID string) (*dto.PlayerSummary, error)
	GetPlayer(ctx context.Context, playerID string) (*dto.PlayerSummary, error)
	ListSkills(ctx context.Context, playerID string) ([]dto.SkillSummary, error)
	UnlockSkill(ctx context.Context, playerID, skillID string) (*dto.SkillActionResult, error)
	EquipSkill(ctx context.Context, playerID, skillID string) (*dto.SkillActionResult, error)
	ActivateSkill(ctx context.Context, playerID, skillID string) (*dto.SkillActionResult, error)
	TickSkillCooldowns(ctx context.Context, playerID string, seconds int) (*dto.SkillActionResult, error)
}

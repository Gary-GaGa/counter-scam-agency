package personnel

import (
	"context"
	"errors"
	"fmt"

	domain "counter-scam-agency/internal/domain/personnel"
	"counter-scam-agency/internal/usecase/dto"
	portin "counter-scam-agency/internal/usecase/port/in"
)

// Compile-time interface guard.
var _ portin.PersonnelUsecase = (*Service)(nil)

var (
	ErrPlayerNotFound   = errors.New("player not found")
	ErrSkillNotFound    = errors.New("skill not found")
	ErrSkillLocked      = errors.New("skill not unlocked")
	ErrSkillNotEquipped = errors.New("skill not equipped")
	ErrSkillOnCooldown  = errors.New("skill on cooldown")
)

// Service orchestrates player/AI skill flows.
type Service struct {
	players domain.PlayerRepository
	catalog map[string]domain.Skill
}

// NewService creates a new personnel service with a skill catalog.
func NewService(players domain.PlayerRepository, catalog []domain.Skill) *Service {
	mapped := make(map[string]domain.Skill, len(catalog))
	for _, skill := range catalog {
		mapped[skill.ID] = skill
	}
	return &Service{players: players, catalog: mapped}
}

// ListSkills lists available skills with player state.
func (s *Service) ListSkills(ctx context.Context, playerID string) ([]dto.SkillSummary, error) {
	player, err := s.players.FindByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("find player: %w", err)
	}
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	results := make([]dto.SkillSummary, 0, len(s.catalog))
	for _, skill := range s.catalog {
		results = append(results, dto.SkillSummary{
			ID:                skill.ID,
			Type:              string(skill.Type),
			Name:              skill.Name,
			Description:       skill.Description,
			CooldownSeconds:   skill.CooldownSeconds,
			RequiredModuleIDs: append([]string{}, skill.RequiredModuleIDs...),
			Unlocked:          player.IsSkillUnlocked(skill.ID),
			Equipped:          player.Partner != nil && player.Partner.HasSkill(skill.ID),
			CooldownRemaining: getCooldown(player, skill.ID),
		})
	}
	return results, nil
}

// UnlockSkill unlocks a skill by ID.
func (s *Service) UnlockSkill(ctx context.Context, playerID, skillID string) (*dto.SkillActionResult, error) {
	player, skill, err := s.loadPlayerAndSkill(ctx, playerID, skillID)
	if err != nil {
		return nil, err
	}
	unlocked := player.UnlockSkill(skill)
	if err := s.players.Save(ctx, player); err != nil {
		return nil, fmt.Errorf("save player: %w", err)
	}
	return skillResult(player, skillID, unlocked, player.Partner != nil && player.Partner.HasSkill(skillID)), nil
}

// EquipSkill equips a learned skill to the AI partner.
func (s *Service) EquipSkill(ctx context.Context, playerID, skillID string) (*dto.SkillActionResult, error) {
	player, skill, err := s.loadPlayerAndSkill(ctx, playerID, skillID)
	if err != nil {
		return nil, err
	}
	if !player.IsSkillUnlocked(skillID) {
		return nil, ErrSkillLocked
	}
	equipped := player.EquipPartnerSkill(skill)
	if err := s.players.Save(ctx, player); err != nil {
		return nil, fmt.Errorf("save player: %w", err)
	}
	return skillResult(player, skillID, player.IsSkillUnlocked(skillID), equipped), nil
}

// ActivateSkill uses a skill if equipped and not on cooldown.
func (s *Service) ActivateSkill(ctx context.Context, playerID, skillID string) (*dto.SkillActionResult, error) {
	player, skill, err := s.loadPlayerAndSkill(ctx, playerID, skillID)
	if err != nil {
		return nil, err
	}
	if !player.IsSkillUnlocked(skillID) {
		return nil, ErrSkillLocked
	}
	if player.Partner == nil || !player.Partner.HasSkill(skillID) {
		return nil, ErrSkillNotEquipped
	}
	if player.Partner.CooldownRemaining(skillID) > 0 {
		return nil, ErrSkillOnCooldown
	}
	if !player.Partner.HasModules(skill.RequiredModuleIDs) {
		return nil, ErrSkillNotEquipped
	}
	if !player.ActivatePartnerSkill(skillID) {
		return nil, ErrSkillOnCooldown
	}
	if err := s.players.Save(ctx, player); err != nil {
		return nil, fmt.Errorf("save player: %w", err)
	}
	return skillResult(player, skillID, true, true), nil
}

// TickSkillCooldowns reduces cooldown timers for all skills.
func (s *Service) TickSkillCooldowns(ctx context.Context, playerID string, seconds int) (*dto.SkillActionResult, error) {
	player, err := s.players.FindByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("find player: %w", err)
	}
	if player == nil {
		return nil, ErrPlayerNotFound
	}
	if player.Partner != nil {
		player.Partner.TickCooldown(seconds)
	}
	if err := s.players.Save(ctx, player); err != nil {
		return nil, fmt.Errorf("save player: %w", err)
	}
	return &dto.SkillActionResult{PlayerID: player.ID}, nil
}

func (s *Service) loadPlayerAndSkill(ctx context.Context, playerID, skillID string) (*domain.Player, domain.Skill, error) {
	player, err := s.players.FindByID(ctx, playerID)
	if err != nil {
		return nil, domain.Skill{}, fmt.Errorf("find player: %w", err)
	}
	if player == nil {
		return nil, domain.Skill{}, ErrPlayerNotFound
	}
	skill, ok := s.catalog[skillID]
	if !ok {
		return nil, domain.Skill{}, ErrSkillNotFound
	}
	return player, skill, nil
}

func skillResult(player *domain.Player, skillID string, unlocked, equipped bool) *dto.SkillActionResult {
	return &dto.SkillActionResult{
		PlayerID:          player.ID,
		SkillID:           skillID,
		Unlocked:          unlocked,
		Equipped:          equipped,
		CooldownRemaining: getCooldown(player, skillID),
	}
}

func getCooldown(player *domain.Player, skillID string) int {
	if player == nil || player.Partner == nil {
		return 0
	}
	return player.Partner.CooldownRemaining(skillID)
}

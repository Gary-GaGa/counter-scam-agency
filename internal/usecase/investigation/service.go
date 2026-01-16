package investigation

import (
	"context"
	"counter-scam-agency/internal/domain/intelligence"
	"counter-scam-agency/internal/domain/operation"
	"counter-scam-agency/internal/domain/personnel"
	"errors"
	"fmt"
)

var (
	ErrInvestigationNotFinished = errors.New("investigation not finished")
	ErrInvestigationNotFound    = errors.New("investigation not found")
	ErrMissionNotFound          = errors.New("mission not found")
	ErrPlayerNotFound           = errors.New("player not found")
)

// Service orchestrates investigation flows.
type Service struct {
	missions       intelligence.MissionRepository
	investigations operation.InvestigationRepository
	players        personnel.PlayerRepository
}

// NewService creates a new investigation service.
func NewService(
	missions intelligence.MissionRepository,
	investigations operation.InvestigationRepository,
	players personnel.PlayerRepository,
) *Service {
	return &Service{
		missions:       missions,
		investigations: investigations,
		players:        players,
	}
}

// CompleteResult summarizes the investigation completion outcome.
type CompleteResult struct {
	InvestigationID  string
	PlayerID         string
	MissionID        string
	Success          bool
	ReputationGained int
}

// CompleteInvestigation finalizes an investigation and applies reputation gains.
func (s *Service) CompleteInvestigation(ctx context.Context, investigationID string) (*CompleteResult, error) {
	inv, err := s.investigations.FindByID(ctx, investigationID)
	if err != nil {
		return nil, fmt.Errorf("find investigation: %w", err)
	}
	if inv == nil {
		return nil, ErrInvestigationNotFound
	}
	if inv.Status == operation.InvestigationStatusActive {
		return nil, ErrInvestigationNotFinished
	}

	mission, err := s.missions.FindByID(ctx, inv.MissionID)
	if err != nil {
		return nil, fmt.Errorf("find mission: %w", err)
	}
	if mission == nil {
		return nil, ErrMissionNotFound
	}

	player, err := s.players.FindByID(ctx, inv.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("find player: %w", err)
	}
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	success := inv.Status == operation.InvestigationStatusCompleted
	gain := mission.ReputationGain(success)
	if gain > 0 {
		player.AddReputation(gain)
		if err := s.players.Save(ctx, player); err != nil {
			return nil, fmt.Errorf("save player: %w", err)
		}
	}

	return &CompleteResult{
		InvestigationID:  inv.ID,
		PlayerID:         inv.PlayerID,
		MissionID:        inv.MissionID,
		Success:          success,
		ReputationGained: gain,
	}, nil
}

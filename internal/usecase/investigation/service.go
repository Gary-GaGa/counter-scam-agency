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
	ErrInvestigationNotActive   = errors.New("investigation not active")
	ErrMissionNotFound          = errors.New("mission not found")
	ErrPlayerNotFound           = errors.New("player not found")
	ErrNodeNotFound             = errors.New("node not found")
	ErrOptionNotFound           = errors.New("option not found")
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

// NodeProgressResult summarizes the node transition result.
type NodeProgressResult struct {
	InvestigationID string
	MissionID       string
	NodeID          string
	OptionID        string
	NextNodeID      string
	Status          operation.InvestigationStatus
}

// AdvanceNode applies a player option, records evidence, and updates investigation status.
func (s *Service) AdvanceNode(ctx context.Context, investigationID, nodeID, optionID string) (*NodeProgressResult, error) {
	inv, err := s.investigations.FindByID(ctx, investigationID)
	if err != nil {
		return nil, fmt.Errorf("find investigation: %w", err)
	}
	if inv == nil {
		return nil, ErrInvestigationNotFound
	}
	if inv.Status != operation.InvestigationStatusActive {
		return nil, ErrInvestigationNotActive
	}

	mission, err := s.missions.FindByID(ctx, inv.MissionID)
	if err != nil {
		return nil, fmt.Errorf("find mission: %w", err)
	}
	if mission == nil {
		return nil, ErrMissionNotFound
	}

	node := mission.GetNode(nodeID)
	if node == nil {
		return nil, ErrNodeNotFound
	}

	var selected *intelligence.NarrativeOption
	for i := range node.Options {
		if node.Options[i].ID == optionID {
			selected = &node.Options[i]
			break
		}
	}
	if selected == nil {
		return nil, ErrOptionNotFound
	}

	for _, evidenceID := range selected.EvidenceIDs {
		inv.CollectEvidence(evidenceID)
	}

	inv.RecordDecision(nodeID, optionID, selected.NextNodeID)
	if selected.LeadsToEnd {
		if selected.SuccessEnd {
			inv.Complete()
		} else {
			inv.Fail()
		}
	}

	if err := s.investigations.Save(ctx, inv); err != nil {
		return nil, fmt.Errorf("save investigation: %w", err)
	}

	return &NodeProgressResult{
		InvestigationID: inv.ID,
		MissionID:       inv.MissionID,
		NodeID:          nodeID,
		OptionID:        optionID,
		NextNodeID:      selected.NextNodeID,
		Status:          inv.Status,
	}, nil
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

package investigation

import (
	"context"
	"counter-scam-agency/internal/domain/intelligence"
	"counter-scam-agency/internal/domain/operation"
	"counter-scam-agency/internal/domain/personnel"
	"counter-scam-agency/internal/usecase/dto"
	"errors"
	"fmt"
)

var (
	ErrInvestigationNotFinished = errors.New("investigation not finished")
	ErrInvestigationNotFound    = errors.New("investigation not found")
	ErrInvestigationNotActive   = errors.New("investigation not active")
	ErrInvestigationExists      = errors.New("investigation already exists")
	ErrMissionNotFound          = errors.New("mission not found")
	ErrPlayerNotFound           = errors.New("player not found")
	ErrNodeNotFound             = errors.New("node not found")
	ErrOptionNotFound           = errors.New("option not found")
	ErrStartNodeNotFound        = errors.New("start node not found")
	ErrEvidenceNotFound         = errors.New("evidence not found")
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

// ListMissions returns available missions for selection.
func (s *Service) ListMissions(ctx context.Context) ([]dto.MissionSummary, error) {
	missions, err := s.missions.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("find missions: %w", err)
	}
	results := make([]dto.MissionSummary, 0, len(missions))
	for _, mission := range missions {
		if mission == nil {
			continue
		}
		results = append(results, dto.MissionSummary{
			ID:               mission.ID,
			Title:            mission.Title,
			Description:      mission.Description,
			Type:             string(mission.Type),
			Difficulty:       mission.Difficulty,
			ReputationWeight: mission.ReputationWeight,
		})
	}
	return results, nil
}

// GetMission returns mission details for UI rendering.
func (s *Service) GetMission(ctx context.Context, missionID string) (*dto.MissionDetail, error) {
	mission, err := s.missions.FindByID(ctx, missionID)
	if err != nil {
		return nil, fmt.Errorf("find mission: %w", err)
	}
	if mission == nil {
		return nil, ErrMissionNotFound
	}

	result := &dto.MissionDetail{
		ID:               mission.ID,
		Title:            mission.Title,
		Description:      mission.Description,
		Type:             string(mission.Type),
		Difficulty:       mission.Difficulty,
		ReputationWeight: mission.ReputationWeight,
		VictimProfile:    mapVictimProfile(mission.VictimProfile),
		Nodes:            make([]dto.NarrativeNode, 0, len(mission.Nodes)),
		EvidenceList:     make([]dto.Evidence, 0, len(mission.EvidenceList)),
	}

	for _, node := range mission.Nodes {
		mapped := dto.NarrativeNode{
			ID:         node.ID,
			Title:      node.Title,
			Body:       node.Body,
			IsTerminal: node.IsTerminal,
			Options:    make([]dto.NarrativeOption, 0, len(node.Options)),
		}
		for _, opt := range node.Options {
			mapped.Options = append(mapped.Options, dto.NarrativeOption{
				ID:          opt.ID,
				Label:       opt.Label,
				NextNodeID:  opt.NextNodeID,
				EvidenceIDs: append([]string{}, opt.EvidenceIDs...),
				LeadsToEnd:  opt.LeadsToEnd,
				SuccessEnd:  opt.SuccessEnd,
			})
		}
		result.Nodes = append(result.Nodes, mapped)
	}

	for _, ev := range mission.EvidenceList {
		result.EvidenceList = append(result.EvidenceList, dto.Evidence{
			ID:              ev.ID,
			Description:     ev.Description,
			Type:            string(ev.Type),
			IsContradiction: ev.IsContradiction,
		})
	}

	return result, nil
}

// StartInvestigation creates a new investigation instance and sets the initial node.
func (s *Service) StartInvestigation(ctx context.Context, investigationID, playerID, missionID, startNodeID string) (*dto.InvestigationStartResult, error) {
	if investigationID != "" {
		existing, err := s.investigations.FindByID(ctx, investigationID)
		if err != nil {
			return nil, fmt.Errorf("find investigation: %w", err)
		}
		if existing != nil {
			return nil, ErrInvestigationExists
		}
	}

	player, err := s.players.FindByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("find player: %w", err)
	}
	if player == nil {
		return nil, ErrPlayerNotFound
	}

	mission, err := s.missions.FindByID(ctx, missionID)
	if err != nil {
		return nil, fmt.Errorf("find mission: %w", err)
	}
	if mission == nil {
		return nil, ErrMissionNotFound
	}

	startNode := resolveStartNode(mission, startNodeID)
	if startNode == nil {
		return nil, ErrStartNodeNotFound
	}

	inv := operation.NewInvestigation(investigationID, playerID, missionID)
	inv.SetCurrentNode(startNode.ID)
	if err := s.investigations.Save(ctx, inv); err != nil {
		return nil, fmt.Errorf("save investigation: %w", err)
	}

	return &dto.InvestigationStartResult{
		InvestigationID: inv.ID,
		PlayerID:        inv.PlayerID,
		MissionID:       inv.MissionID,
		CurrentNodeID:   inv.CurrentNodeID,
		Status:          string(inv.Status),
	}, nil
}

// AdvanceNode applies a player option, records evidence, and updates investigation status.
func (s *Service) AdvanceNode(ctx context.Context, investigationID, nodeID, optionID string) (*dto.NodeProgressResult, error) {
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

	return &dto.NodeProgressResult{
		InvestigationID: inv.ID,
		MissionID:       inv.MissionID,
		NodeID:          nodeID,
		OptionID:        optionID,
		NextNodeID:      selected.NextNodeID,
		Status:          string(inv.Status),
	}, nil
}

// SubmitEvidence records evidence submission and returns contradiction info.
func (s *Service) SubmitEvidence(ctx context.Context, investigationID, evidenceID string) (*dto.SubmitEvidenceResult, error) {
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

	evidence := mission.GetEvidence(evidenceID)
	if evidence == nil {
		return nil, ErrEvidenceNotFound
	}

	alreadyCollected := isEvidenceCollected(inv, evidenceID)
	inv.CollectEvidence(evidenceID)
	if !evidence.IsContradiction {
		inv.IncreaseSuspicion(10)
	}
	if err := s.investigations.Save(ctx, inv); err != nil {
		return nil, fmt.Errorf("save investigation: %w", err)
	}

	return &dto.SubmitEvidenceResult{
		InvestigationID:  inv.ID,
		EvidenceID:       evidenceID,
		IsContradiction:  evidence.IsContradiction,
		AlreadyCollected: alreadyCollected,
	}, nil
}

// CompleteInvestigation finalizes an investigation and applies reputation gains.
func (s *Service) CompleteInvestigation(ctx context.Context, investigationID string) (*dto.CompleteResult, error) {
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

	return &dto.CompleteResult{
		InvestigationID:  inv.ID,
		PlayerID:         inv.PlayerID,
		MissionID:        inv.MissionID,
		Success:          success,
		ReputationGained: gain,
	}, nil
}

func resolveStartNode(mission *intelligence.Mission, explicitID string) *intelligence.NarrativeNode {
	if mission == nil {
		return nil
	}
	if explicitID != "" {
		return mission.GetNode(explicitID)
	}
	if node := mission.GetNode("start"); node != nil {
		return node
	}
	if len(mission.Nodes) > 0 {
		return &mission.Nodes[0]
	}
	return nil
}

func isEvidenceCollected(inv *operation.Investigation, evidenceID string) bool {
	if inv == nil {
		return false
	}
	for _, id := range inv.CollectedEvidenceIDs {
		if id == evidenceID {
			return true
		}
	}
	return false
}

func mapVictimProfile(profile *intelligence.VictimProfile) *dto.VictimProfile {
	if profile == nil {
		return nil
	}
	return &dto.VictimProfile{
		Anxiety:   profile.Anxiety,
		Trust:     profile.Trust,
		Urgency:   profile.Urgency,
		Isolation: profile.Isolation,
		RiskScore: profile.RiskScore(),
		RiskLevel: string(profile.RiskLevel()),
	}
}

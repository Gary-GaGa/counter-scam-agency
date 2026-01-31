package in

import (
	"context"

	"counter-scam-agency/internal/usecase/dto"
)

// InvestigationUsecase defines input port for investigation flows.
type InvestigationUsecase interface {
	ListMissions(ctx context.Context) ([]dto.MissionSummary, error)
	GetMission(ctx context.Context, missionID string) (*dto.MissionDetail, error)
	StartInvestigation(ctx context.Context, investigationID, playerID, missionID, startNodeID string) (*dto.InvestigationStartResult, error)
	AdvanceNode(ctx context.Context, investigationID, nodeID, optionID string) (*dto.NodeProgressResult, error)
	SubmitEvidence(ctx context.Context, investigationID, evidenceID string) (*dto.SubmitEvidenceResult, error)
	CompleteInvestigation(ctx context.Context, investigationID string) (*dto.CompleteResult, error)
}

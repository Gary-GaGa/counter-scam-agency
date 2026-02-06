package model

import (
	model "counter-scam-agency/internal/domain/operation"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/operation"
)

// InvestigationDocToModel converts a MongoDB investigation document to a model.
func InvestigationDocToModel(in *operation.InvestigationDocPo) *model.Investigation {
	if in == nil {
		return new(model.Investigation)
	}
	return &model.Investigation{
		ID:                   in.ID,
		PlayerID:             in.PlayerID,
		MissionID:            in.MissionID,
		Status:               model.InvestigationStatus(in.Status),
		CurrentNodeID:        in.CurrentNodeID,
		NodeHistory:          NodeDecisionsToModel(in.NodeHistory),
		CollectedEvidenceIDs: append([]string{}, in.CollectedEvidenceIDs...),
		SuspicionLevel:       in.SuspicionLevel,
	}
}

// NodeDecisionsToModel converts node decisions to model.
func NodeDecisionsToModel(in []operation.NodeDecisionDocPo) []model.NodeDecision {
	out := make([]model.NodeDecision, len(in))
	for i, decision := range in {
		out[i] = NodeDecisionToModel(decision)
	}
	return out
}

// NodeDecisionToModel converts a node decision to model.
func NodeDecisionToModel(in operation.NodeDecisionDocPo) model.NodeDecision {
	return model.NodeDecision{
		NodeID:     in.NodeID,
		OptionID:   in.OptionID,
		NextNodeID: in.NextNodeID,
	}
}

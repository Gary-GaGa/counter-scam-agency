package po

import (
	model "counter-scam-agency/internal/domain/operation"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/operation"
)

// InvestigationDocToPo converts an investigation model to PO.
func InvestigationDocToPo(in *model.Investigation) *operation.InvestigationDocPo {
	if in == nil {
		return new(operation.InvestigationDocPo)
	}
	return &operation.InvestigationDocPo{
		ID:                   in.ID,
		PlayerID:             in.PlayerID,
		MissionID:            in.MissionID,
		Status:               string(in.Status),
		CurrentNodeID:        in.CurrentNodeID,
		NodeHistory:          NodeDecisionsToPo(in.NodeHistory),
		CollectedEvidenceIDs: append([]string{}, in.CollectedEvidenceIDs...),
		SuspicionLevel:       in.SuspicionLevel,
	}
}

// NodeDecisionsToPo converts node decisions to PO.
func NodeDecisionsToPo(in []model.NodeDecision) []operation.NodeDecisionDocPo {
	out := make([]operation.NodeDecisionDocPo, len(in))
	for i, decision := range in {
		out[i] = NodeDecisionToPo(decision)
	}
	return out
}

// NodeDecisionToPo converts a node decision to PO.
func NodeDecisionToPo(in model.NodeDecision) operation.NodeDecisionDocPo {
	return operation.NodeDecisionDocPo{
		NodeID:     in.NodeID,
		OptionID:   in.OptionID,
		NextNodeID: in.NextNodeID,
	}
}

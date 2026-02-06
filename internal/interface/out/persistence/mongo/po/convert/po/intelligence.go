package po

import (
	model "counter-scam-agency/internal/domain/intelligence"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/intelligence"
)

// MissionDocToPo converts a mission model to a MongoDB document.
func MissionDocToPo(in *model.Mission) *intelligence.MissionDocPo {
	if in == nil {
		return new(intelligence.MissionDocPo)
	}

	return &intelligence.MissionDocPo{
		ID:               in.ID,
		Title:            in.Title,
		Description:      in.Description,
		Type:             string(in.Type),
		Difficulty:       in.Difficulty,
		ReputationWeight: in.ReputationWeight,
		VictimProfile:    VictimProfileDocToPo(in.VictimProfile),
		Nodes:            NarrativeNodesToPo(in.Nodes),
		EvidenceList:     EvidenceDocsToPo(in.EvidenceList),
	}
}

// VictimProfileDocToPo converts victim profile to PO.
func VictimProfileDocToPo(in *model.VictimProfile) *intelligence.VictimProfileDocPo {
	if in == nil {
		return nil
	}
	return &intelligence.VictimProfileDocPo{
		Anxiety:   in.Anxiety,
		Trust:     in.Trust,
		Urgency:   in.Urgency,
		Isolation: in.Isolation,
	}
}

// NarrativeNodesToPo converts narrative nodes to PO.
func NarrativeNodesToPo(in []model.NarrativeNode) []intelligence.NarrativeNodeDocPo {
	out := make([]intelligence.NarrativeNodeDocPo, len(in))
	for i, node := range in {
		out[i] = NarrativeNodeToPo(node)
	}
	return out
}

// NarrativeNodeToPo converts a narrative node to PO.
func NarrativeNodeToPo(in model.NarrativeNode) intelligence.NarrativeNodeDocPo {
	return intelligence.NarrativeNodeDocPo{
		ID:         in.ID,
		Title:      in.Title,
		Body:       in.Body,
		Options:    NarrativeOptionsToPo(in.Options),
		IsTerminal: in.IsTerminal,
	}
}

// NarrativeOptionsToPo converts narrative options to PO.
func NarrativeOptionsToPo(in []model.NarrativeOption) []intelligence.NarrativeOptionDocPo {
	out := make([]intelligence.NarrativeOptionDocPo, len(in))
	for i, opt := range in {
		out[i] = NarrativeOptionToPo(opt)
	}
	return out
}

// NarrativeOptionToPo converts a narrative option to PO.
func NarrativeOptionToPo(in model.NarrativeOption) intelligence.NarrativeOptionDocPo {
	return intelligence.NarrativeOptionDocPo{
		ID:          in.ID,
		Label:       in.Label,
		NextNodeID:  in.NextNodeID,
		EvidenceIDs: append([]string{}, in.EvidenceIDs...),
		LeadsToEnd:  in.LeadsToEnd,
		SuccessEnd:  in.SuccessEnd,
	}
}

// EvidenceDocsToPo converts evidence list to PO.
func EvidenceDocsToPo(in []model.Evidence) []intelligence.EvidenceDocPo {
	out := make([]intelligence.EvidenceDocPo, len(in))
	for i, ev := range in {
		out[i] = EvidenceDocToPo(ev)
	}
	return out
}

// EvidenceDocToPo converts evidence to PO.
func EvidenceDocToPo(in model.Evidence) intelligence.EvidenceDocPo {
	return intelligence.EvidenceDocPo{
		ID:              in.ID,
		Description:     in.Description,
		Type:            string(in.Type),
		IsContradiction: in.IsContradiction,
	}
}

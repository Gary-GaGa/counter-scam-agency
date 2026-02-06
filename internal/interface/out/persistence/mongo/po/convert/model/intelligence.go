package model

import (
	model "counter-scam-agency/internal/domain/intelligence"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/intelligence"
)

// MissionDocToModel converts a MongoDB mission document to a model.
func MissionDocToModel(in *intelligence.MissionDocPo) *model.Mission {
	if in == nil {
		return new(model.Mission)
	}

	return &model.Mission{
		ID:               in.ID,
		Title:            in.Title,
		Description:      in.Description,
		Type:             model.ScamType(in.Type),
		Difficulty:       in.Difficulty,
		ReputationWeight: in.ReputationWeight,
		VictimProfile:    VictimProfileDocToModel(in.VictimProfile),
		Nodes:            NarrativeNodesToModel(in.Nodes),
		EvidenceList:     EvidenceDocsToModel(in.EvidenceList),
	}
}

// VictimProfileDocToModel converts victim profile to model.
func VictimProfileDocToModel(in *intelligence.VictimProfileDocPo) *model.VictimProfile {
	if in == nil {
		return nil
	}
	return &model.VictimProfile{
		Anxiety:   in.Anxiety,
		Trust:     in.Trust,
		Urgency:   in.Urgency,
		Isolation: in.Isolation,
	}
}

// NarrativeNodesToModel converts narrative nodes to model.
func NarrativeNodesToModel(in []intelligence.NarrativeNodeDocPo) []model.NarrativeNode {
	out := make([]model.NarrativeNode, len(in))
	for i, node := range in {
		out[i] = NarrativeNodeToModel(node)
	}
	return out
}

// NarrativeNodeToModel converts a narrative node to model.
func NarrativeNodeToModel(in intelligence.NarrativeNodeDocPo) model.NarrativeNode {
	return model.NarrativeNode{
		ID:         in.ID,
		Title:      in.Title,
		Body:       in.Body,
		Options:    NarrativeOptionsToModel(in.Options),
		IsTerminal: in.IsTerminal,
	}
}

// NarrativeOptionsToModel converts narrative options to model.
func NarrativeOptionsToModel(in []intelligence.NarrativeOptionDocPo) []model.NarrativeOption {
	out := make([]model.NarrativeOption, len(in))
	for i, opt := range in {
		out[i] = NarrativeOptionToModel(opt)
	}
	return out
}

// NarrativeOptionToModel converts a narrative option to model.
func NarrativeOptionToModel(in intelligence.NarrativeOptionDocPo) model.NarrativeOption {
	return model.NarrativeOption{
		ID:          in.ID,
		Label:       in.Label,
		NextNodeID:  in.NextNodeID,
		EvidenceIDs: append([]string{}, in.EvidenceIDs...),
		LeadsToEnd:  in.LeadsToEnd,
		SuccessEnd:  in.SuccessEnd,
	}
}

// EvidenceDocsToModel converts evidence list to model.
func EvidenceDocsToModel(in []intelligence.EvidenceDocPo) []model.Evidence {
	out := make([]model.Evidence, len(in))
	for i, ev := range in {
		out[i] = EvidenceDocToModel(ev)
	}
	return out
}

// EvidenceDocToModel converts evidence to model.
func EvidenceDocToModel(in intelligence.EvidenceDocPo) model.Evidence {
	return model.Evidence{
		ID:              in.ID,
		Description:     in.Description,
		Type:            model.EvidenceType(in.Type),
		IsContradiction: in.IsContradiction,
	}
}

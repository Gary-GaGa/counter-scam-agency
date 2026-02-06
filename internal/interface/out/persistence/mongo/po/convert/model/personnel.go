package model

import (
	model "counter-scam-agency/internal/domain/personnel"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/personnel"
)

// PlayerDocToModel converts a MongoDB player document to a model.
func PlayerDocToModel(in *personnel.PlayerDocPo) *model.Player {
	if in == nil {
		return new(model.Player)
	}
	return &model.Player{
		ID:                in.ID,
		Stats:             StatsToModel(in.Stats),
		Partner:           AIPartnerToModel(in.Partner),
		Reputation:        in.Reputation,
		UnlockedModuleIDs: sliceToKeySet(in.UnlockedModuleIDs),
		UnlockedSkillIDs:  sliceToKeySet(in.UnlockedSkillIDs),
	}
}

// AIPartnerToModel converts AI partner to model.
func AIPartnerToModel(in *personnel.AIPartnerDocPo) *model.AIPartner {
	if in == nil {
		return nil
	}
	return &model.AIPartner{
		Loadout:        ModulesToModel(in.Loadout),
		Skills:         SkillsToModel(in.Skills),
		SkillCooldowns: cloneCooldowns(in.SkillCooldowns),
		Personality:    model.PersonalityType(in.Personality),
	}
}

// ModulesToModel converts modules to model.
func ModulesToModel(in []personnel.ModuleDocPo) []model.Module {
	out := make([]model.Module, len(in))
	for i, mod := range in {
		out[i] = ModuleToModel(mod)
	}
	return out
}

// ModuleToModel converts a module to model.
func ModuleToModel(in personnel.ModuleDocPo) model.Module {
	return model.Module{
		ID:                 in.ID,
		Type:               model.ModuleType(in.Type),
		Name:               in.Name,
		Description:        in.Description,
		Level:              in.Level,
		ReputationRequired: in.ReputationRequired,
		StatBonus:          StatsToModel(in.StatBonus),
	}
}

// SkillsToModel converts skills to model.
func SkillsToModel(in []personnel.SkillDocPo) []model.Skill {
	out := make([]model.Skill, len(in))
	for i, skill := range in {
		out[i] = SkillToModel(skill)
	}
	return out
}

// SkillToModel converts a skill to model.
func SkillToModel(in personnel.SkillDocPo) model.Skill {
	return model.Skill{
		ID:                 in.ID,
		Type:               model.SkillType(in.Type),
		Name:               in.Name,
		Description:        in.Description,
		CooldownSeconds:    in.CooldownSeconds,
		ReputationRequired: in.ReputationRequired,
		RequiredModuleIDs:  append([]string{}, in.RequiredModuleIDs...),
	}
}

// StatsToModel converts stats to model.
func StatsToModel(in personnel.StatsDocPo) model.Stats {
	return model.Stats{
		Logic:      in.Logic,
		Tech:       in.Tech,
		Charisma:   in.Charisma,
		Resilience: in.Resilience,
	}
}

func sliceToKeySet(in []string) map[string]struct{} {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(in))
	for _, value := range in {
		out[value] = struct{}{}
	}
	return out
}

func cloneCooldowns(in map[string]int) map[string]int {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]int, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

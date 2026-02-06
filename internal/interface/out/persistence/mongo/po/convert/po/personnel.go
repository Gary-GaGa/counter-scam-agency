package po

import (
	model "counter-scam-agency/internal/domain/personnel"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/personnel"
)

// PlayerDocToPo converts a player model to PO.
func PlayerDocToPo(in *model.Player) *personnel.PlayerDocPo {
	if in == nil {
		return new(personnel.PlayerDocPo)
	}
	return &personnel.PlayerDocPo{
		ID:                in.ID,
		Stats:             StatsToPo(in.Stats),
		Partner:           AIPartnerToPo(in.Partner),
		Reputation:        in.Reputation,
		UnlockedModuleIDs: mapKeySetToSlice(in.UnlockedModuleIDs),
		UnlockedSkillIDs:  mapKeySetToSlice(in.UnlockedSkillIDs),
	}
}

// AIPartnerToPo converts AI partner to PO.
func AIPartnerToPo(in *model.AIPartner) *personnel.AIPartnerDocPo {
	if in == nil {
		return nil
	}
	return &personnel.AIPartnerDocPo{
		Loadout:        ModulesToPo(in.Loadout),
		Skills:         SkillsToPo(in.Skills),
		SkillCooldowns: cloneCooldowns(in.SkillCooldowns),
		Personality:    string(in.Personality),
	}
}

// ModulesToPo converts modules to PO.
func ModulesToPo(in []model.Module) []personnel.ModuleDocPo {
	out := make([]personnel.ModuleDocPo, len(in))
	for i, mod := range in {
		out[i] = ModuleToPo(mod)
	}
	return out
}

// ModuleToPo converts a module to PO.
func ModuleToPo(in model.Module) personnel.ModuleDocPo {
	return personnel.ModuleDocPo{
		ID:                 in.ID,
		Type:               string(in.Type),
		Name:               in.Name,
		Description:        in.Description,
		Level:              in.Level,
		ReputationRequired: in.ReputationRequired,
		StatBonus:          StatsToPo(in.StatBonus),
	}
}

// SkillsToPo converts skills to PO.
func SkillsToPo(in []model.Skill) []personnel.SkillDocPo {
	out := make([]personnel.SkillDocPo, len(in))
	for i, skill := range in {
		out[i] = SkillToPo(skill)
	}
	return out
}

// SkillToPo converts a skill to PO.
func SkillToPo(in model.Skill) personnel.SkillDocPo {
	return personnel.SkillDocPo{
		ID:                 in.ID,
		Type:               string(in.Type),
		Name:               in.Name,
		Description:        in.Description,
		CooldownSeconds:    in.CooldownSeconds,
		ReputationRequired: in.ReputationRequired,
		RequiredModuleIDs:  append([]string{}, in.RequiredModuleIDs...),
	}
}

// StatsToPo converts stats to PO.
func StatsToPo(in model.Stats) personnel.StatsDocPo {
	return personnel.StatsDocPo{
		Logic:      in.Logic,
		Tech:       in.Tech,
		Charisma:   in.Charisma,
		Resilience: in.Resilience,
	}
}

func mapKeySetToSlice(in map[string]struct{}) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	for key := range in {
		out = append(out, key)
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

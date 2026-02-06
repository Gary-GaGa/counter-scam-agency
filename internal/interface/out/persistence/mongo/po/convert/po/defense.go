package po

import (
	model "counter-scam-agency/internal/domain/defense"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/defense"
)

// BaseDocToPo converts a domain Base to a BaseDocPo.
func BaseDocToPo(in *model.Base) *defense.BaseDocPo {
	if in == nil {
		return new(defense.BaseDocPo)
	}

	return &defense.BaseDocPo{
		ID:            in.ID,
		OwnerID:       in.OwnerID,
		SecurityLevel: in.SecurityLevel,
		FacilitySlots: in.FacilitySlots,
		Facilities:    FacilityDocsToPo(in.Facilities),
		Upgrades:      UpgradeDocsToPo(in.Upgrades),
	}
}

// FacilityDocsToPo converts a slice of domain Facility to FacilityDocPo.
func FacilityDocsToPo(in []model.Facility) []defense.FacilityDocPo {
	list := make([]defense.FacilityDocPo, len(in))

	for i, item := range in {
		list[i] = FacilityDocToPo(item)
	}

	return list
}

// FacilityDocToPo converts a single domain Facility to FacilityDocPo.
func FacilityDocToPo(in model.Facility) defense.FacilityDocPo {
	return defense.FacilityDocPo{
		ID:          in.ID,
		Type:        string(in.Type),
		Name:        in.Name,
		Level:       in.Level,
		MaxLevel:    in.MaxLevel,
		Description: in.Description,
	}
}

// UpgradeDocsToPo converts a slice of domain Upgrade to UpgradeDocPo.
func UpgradeDocsToPo(in []model.Upgrade) []defense.UpgradeDocPo {
	list := make([]defense.UpgradeDocPo, len(in))

	for i, item := range in {
		list[i] = UpgradeDocToPo(item)
	}

	return list
}

// UpgradeDocToPo converts a single domain Upgrade to UpgradeDocPo.
func UpgradeDocToPo(in model.Upgrade) defense.UpgradeDocPo {
	return defense.UpgradeDocPo{
		ID:          in.ID,
		Name:        in.Name,
		Level:       in.Level,
		MaxLevel:    in.MaxLevel,
		Description: in.Description,
	}
}

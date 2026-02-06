package model

import (
	model "counter-scam-agency/internal/domain/defense"
	"counter-scam-agency/internal/interface/out/persistence/mongo/po/defense"
)

// BaseDocToModel -
func BaseDocToModel(in *defense.BaseDocPo) *model.Base {
	if in == nil {
		return new(model.Base)
	}

	return &model.Base{
		ID:            in.ID,
		OwnerID:       in.OwnerID,
		SecurityLevel: in.SecurityLevel,
		FacilitySlots: in.FacilitySlots,
		Facilities:    FacilityDocsToModel(in.Facilities),
		Upgrades:      UpgradeDocsToModel(in.Upgrades),
	}
}

// FacilityDocsToModel -
func FacilityDocsToModel(in []defense.FacilityDocPo) []model.Facility {
	list := make([]model.Facility, len(in))

	for i, item := range in {
		list[i] = FacilityDocToModel(item)
	}

	return list
}

// FacilityDocToModel -
func FacilityDocToModel(in defense.FacilityDocPo) model.Facility {
	return model.Facility{
		ID:          in.ID,
		Type:        model.FacilityType(in.Type),
		Name:        in.Name,
		Level:       in.Level,
		MaxLevel:    in.MaxLevel,
		Description: in.Description,
	}
}

// UpgradeDocsToModel -
func UpgradeDocsToModel(in []defense.UpgradeDocPo) []model.Upgrade {
	list := make([]model.Upgrade, len(in))

	for i, item := range in {
		list[i] = UpgradeDocToModel(item)
	}

	return list
}

// UpgradeDocToModel -
func UpgradeDocToModel(in defense.UpgradeDocPo) model.Upgrade {
	return model.Upgrade{
		ID:          in.ID,
		Name:        in.Name,
		Level:       in.Level,
		MaxLevel:    in.MaxLevel,
		Description: in.Description,
	}
}

package intelligence

// RiskLevel describes the overall vulnerability level.
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "Low"
	RiskLevelMedium RiskLevel = "Medium"
	RiskLevelHigh   RiskLevel = "High"
)

// VictimProfile represents the psychological profile of a victim.
type VictimProfile struct {
	Anxiety   int // 0-100, higher means more anxious
	Trust     int // 0-100, higher means more trusting
	Urgency   int // 0-100, higher means more impulsive
	Isolation int // 0-100, higher means more isolated
}

// RiskScore returns an average risk score (0-100).
func (p VictimProfile) RiskScore() int {
	total := clampScore(p.Anxiety) + clampScore(p.Trust) + clampScore(p.Urgency) + clampScore(p.Isolation)
	return total / 4
}

// RiskLevel returns the overall risk level based on the score.
func (p VictimProfile) RiskLevel() RiskLevel {
	score := p.RiskScore()
	switch {
	case score >= 67:
		return RiskLevelHigh
	case score >= 34:
		return RiskLevelMedium
	default:
		return RiskLevelLow
	}
}

func clampScore(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

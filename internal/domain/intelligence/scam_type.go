package intelligence

// ScamType defines the category of the scam.
type ScamType string

const (
	ScamTypePhishing      ScamType = "Phishing"
	ScamTypeInvestment    ScamType = "Investment"
	ScamTypeRomance       ScamType = "Romance"
	ScamTypeImpersonation ScamType = "Impersonation"
)

// String returns the string representation of the ScamType.
func (s ScamType) String() string {
	return string(s)
}

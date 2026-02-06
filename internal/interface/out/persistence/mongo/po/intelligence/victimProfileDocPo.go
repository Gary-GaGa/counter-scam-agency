package intelligence

// VictimProfileDocPo stores victim profile data.
type VictimProfileDocPo struct {
	Anxiety   int `bson:"anxiety"`
	Trust     int `bson:"trust"`
	Urgency   int `bson:"urgency"`
	Isolation int `bson:"isolation"`
}

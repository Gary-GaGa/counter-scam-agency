package personnel

// Stats represents the four core attributes of the player.
type Stats struct {
	Logic      int // Influences contradiction detection
	Tech       int // Influences tool usage success rate
	Charisma   int // Influences NPC trust and negotiation
	Resilience int // Influences SAN value and stress resistance
}

// Add returns a new Stats object with the sum of two stats.
func (s Stats) Add(other Stats) Stats {
	return Stats{
		Logic:      s.Logic + other.Logic,
		Tech:       s.Tech + other.Tech,
		Charisma:   s.Charisma + other.Charisma,
		Resilience: s.Resilience + other.Resilience,
	}
}

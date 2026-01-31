package personnel

import "context"

// PlayerRepository defines the persistence operations for Player aggregates.
type PlayerRepository interface {
	// Save persists the player's state.
	Save(ctx context.Context, player *Player) error

	// FindByID retrieves a player by their ID.
	FindByID(ctx context.Context, id string) (*Player, error)
}

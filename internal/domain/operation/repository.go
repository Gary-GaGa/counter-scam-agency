package operation

import "context"

// InvestigationRepository defines the persistence operations for Investigation aggregates.
type InvestigationRepository interface {
	// Save persists the investigation state.
	Save(ctx context.Context, inv *Investigation) error

	// FindByID retrieves an investigation by its ID.
	FindByID(ctx context.Context, id string) (*Investigation, error)

	// FindByPlayerID retrieves all investigations for a specific player.
	FindByPlayerID(ctx context.Context, playerID string) ([]*Investigation, error)
}

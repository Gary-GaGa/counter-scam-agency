package intelligence

import "context"

// MissionRepository defines the persistence operations for Mission aggregates.
type MissionRepository interface {
	// FindByID retrieves a mission by its ID.
	FindByID(ctx context.Context, id string) (*Mission, error)

	// FindAll retrieves all available missions.
	FindAll(ctx context.Context) ([]*Mission, error)

	// Save persists a mission (useful for seeding or admin tools).
	Save(ctx context.Context, mission *Mission) error
}

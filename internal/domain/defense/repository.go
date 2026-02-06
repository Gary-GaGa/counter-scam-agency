package defense

import "context"

// BaseRepository -
type BaseRepository interface {
	Save(ctx context.Context, base *Base) error
	FindByID(ctx context.Context, id string) (*Base, error)
	FindByOwnerID(ctx context.Context, ownerID string) ([]*Base, error)
}

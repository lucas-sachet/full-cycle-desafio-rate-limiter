package limiter

import "context"

type PersistenceStrategy interface {
	Increment(ctx context.Context, key string, expirationSeconds int) (int, error)
	IsBlocked(ctx context.Context, key string) (bool, error)
	Block(ctx context.Context, key string, durationSeconds int) error
}

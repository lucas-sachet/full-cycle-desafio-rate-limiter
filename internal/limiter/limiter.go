package limiter

import (
	"context"
	"fmt"
)

type RateLimiter struct {
	strategy             PersistenceStrategy
	defaultIPLimit       int
	defaultTokenLimit    int
	blockDurationSeconds int
}

func NewRateLimiter(strategy PersistenceStrategy, ipLimit, tokenLimit, blockDuration int) *RateLimiter {
	return &RateLimiter{
		strategy:             strategy,
		defaultIPLimit:       ipLimit,
		defaultTokenLimit:    tokenLimit,
		blockDurationSeconds: blockDuration,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, ip, token string) (bool, error) {
	// Precedence: Token > IP
	var key string
	var limit int

	if token != "" {
		key = fmt.Sprintf("token:%s", token)
		limit = rl.defaultTokenLimit
		// In a real scenario, you might want to fetch a specific limit for this token from a DB.
		// For this challenge, we'll use a default token limit or just follow the requirement
		// that token overrides IP.
	} else {
		key = fmt.Sprintf("ip:%s", ip)
		limit = rl.defaultIPLimit
	}

	// Check if already blocked
	blocked, err := rl.strategy.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	// Increment and check limit
	// We use a 1-second window for "requests per second"
	count, err := rl.strategy.Increment(ctx, key, 1)
	if err != nil {
		return false, err
	}

	if count > limit {
		// Block the key
		err = rl.strategy.Block(ctx, key, rl.blockDurationSeconds)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

package limiter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPersistenceStrategy struct {
	mock.Mock
}

func (m *MockPersistenceStrategy) Increment(ctx context.Context, key string, expirationSeconds int) (int, error) {
	args := m.Called(ctx, key, expirationSeconds)
	return args.Int(0), args.Error(1)
}

func (m *MockPersistenceStrategy) IsBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockPersistenceStrategy) Block(ctx context.Context, key string, durationSeconds int) error {
	args := m.Called(ctx, key, durationSeconds)
	return args.Error(0)
}

func TestRateLimiter_Allow(t *testing.T) {
	ctx := context.Background()

	t.Run("Allow IP within limit", func(t *testing.T) {
		mockStrategy := new(MockPersistenceStrategy)
		rl := NewRateLimiter(mockStrategy, 10, 100, 300)

		mockStrategy.On("IsBlocked", ctx, "ip:127.0.0.1").Return(false, nil)
		mockStrategy.On("Increment", ctx, "ip:127.0.0.1", 1).Return(5, nil)

		allowed, err := rl.Allow(ctx, "127.0.0.1", "")
		assert.NoError(t, err)
		assert.True(t, allowed)
		mockStrategy.AssertExpectations(t)
	})

	t.Run("Block IP exceeding limit", func(t *testing.T) {
		mockStrategy := new(MockPersistenceStrategy)
		rl := NewRateLimiter(mockStrategy, 10, 100, 300)

		mockStrategy.On("IsBlocked", ctx, "ip:127.0.0.1").Return(false, nil)
		mockStrategy.On("Increment", ctx, "ip:127.0.0.1", 1).Return(11, nil)
		mockStrategy.On("Block", ctx, "ip:127.0.0.1", 300).Return(nil)

		allowed, err := rl.Allow(ctx, "127.0.0.1", "")
		assert.NoError(t, err)
		assert.False(t, allowed)
		mockStrategy.AssertExpectations(t)
	})

	t.Run("Allow Token within limit", func(t *testing.T) {
		mockStrategy := new(MockPersistenceStrategy)
		rl := NewRateLimiter(mockStrategy, 10, 100, 300)

		mockStrategy.On("IsBlocked", ctx, "token:my-token").Return(false, nil)
		mockStrategy.On("Increment", ctx, "token:my-token", 1).Return(50, nil)

		allowed, err := rl.Allow(ctx, "127.0.0.1", "my-token")
		assert.NoError(t, err)
		assert.True(t, allowed)
		mockStrategy.AssertExpectations(t)
	})

	t.Run("Precedence: Token overrides IP limit", func(t *testing.T) {
		mockStrategy := new(MockPersistenceStrategy)
		// IP limit is 2, Token limit is 5
		rl := NewRateLimiter(mockStrategy, 2, 5, 300)

		// Even if IP would be blocked if it was the only key, token takes precedence
		mockStrategy.On("IsBlocked", ctx, "token:my-token").Return(false, nil)
		mockStrategy.On("Increment", ctx, "token:my-token", 1).Return(3, nil)

		allowed, err := rl.Allow(ctx, "127.0.0.1", "my-token")
		assert.NoError(t, err)
		assert.True(t, allowed) // 3 < 5 (token limit), so allowed
		mockStrategy.AssertExpectations(t)
	})

	t.Run("Is already blocked", func(t *testing.T) {
		mockStrategy := new(MockPersistenceStrategy)
		rl := NewRateLimiter(mockStrategy, 10, 100, 300)

		mockStrategy.On("IsBlocked", ctx, "ip:127.0.0.1").Return(true, nil)

		allowed, err := rl.Allow(ctx, "127.0.0.1", "")
		assert.NoError(t, err)
		assert.False(t, allowed)
		mockStrategy.AssertExpectations(t)
	})
}

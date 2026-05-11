package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStrategy struct {
	client *redis.Client
}

func NewRedisStrategy(host, port, password string, db int) *RedisStrategy {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	return &RedisStrategy{client: rdb}
}

func (r *RedisStrategy) Increment(ctx context.Context, key string, expirationSeconds int) (int, error) {
	pipe := r.client.TxPipeline()
	count := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Duration(expirationSeconds)*time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return int(count.Val()), nil
}

func (r *RedisStrategy) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	val, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *RedisStrategy) Block(ctx context.Context, key string, durationSeconds int) error {
	blockKey := fmt.Sprintf("block:%s", key)
	return r.client.Set(ctx, blockKey, "1", time.Duration(durationSeconds)*time.Second).Err()
}
